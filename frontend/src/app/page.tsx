"use client";

import { PaginationControl } from "@/app/components/PaginationControl";
import { TodoItem } from "@/app/components/TodoItem";
import { TodoModal } from "@/app/components/TodoModal";
import { Button } from "@/app/components/ui/button";
import { api } from "@/lib/api";
import { ListTodoResponse, Todo } from "@/types";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Plus } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { useState } from "react";
import { toast } from "sonner";

export default function Home() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const pageParam = searchParams.get("page");
  const currentPage = pageParam ? parseInt(pageParam, 10) : 1;

  const queryClient = useQueryClient();
  const { data, isLoading, error } = useQuery<ListTodoResponse>({
    queryKey: ["todos", currentPage],
    queryFn: async () => {
      const res = await api.get<ListTodoResponse>("/todo", {
        params: { page: currentPage },
      });
      return res.data;
    },
  });

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null);
  const [modalKey, setModalKey] = useState(0);

  const toastHandler = (message: string, type: "success" | "error") => {
    if (type === "success") {
      toast.success(message, {
        duration: 6000,
        position: "bottom-center",
      });
    } else {
      toast.error(message, {
        duration: 6000,
        position: "top-center",
      });
    }
  };
  // Mutations
  const createMutation = useMutation({
    mutationFn: async (data: { title: string; description: string }) => {
      const res = await api.post("/todo", data);
      return res.data;
    },
    onSuccess: () => {
      // Invalidate current page to reflect changes (items might shift)
      queryClient.invalidateQueries({ queryKey: ["todos", currentPage] });
      setIsModalOpen(false);
      toastHandler("Todo created successfully", "success");
    },
    onError: () => {
      toastHandler("Failed to create todo", "error");
    },
  });

  const updateMutation = useMutation({
    mutationFn: async ({
      id,
      data,
    }: {
      id: number;
      data: { title: string; description: string };
    }) => {
      const res = await api.patch(`/todo/${id}`, data);
      return res.data;
    },
    onSuccess: (updatedTodo) => {
      queryClient.invalidateQueries({ queryKey: ["todos", currentPage] });
      // Update selected todo if currently viewing it to reflect changes (e.g. updated_at if displayed)
      setSelectedTodo(updatedTodo);
      // Modal stays open
      toastHandler("Todo updated successfully", "success");
    },
    onError: () => {
      toastHandler("Failed to update todo", "error");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/todo/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos", currentPage] });
      toastHandler("Todo deleted successfully", "success");
    },
    onError: () => {
      toastHandler("Failed to delete todo", "error");
    },
  });

  const handleOpenCreate = () => {
    setSelectedTodo(null);
    setModalKey((prev) => prev + 1);
    setIsModalOpen(true);
  };

  const handleOpenEdit = (todo: Todo) => {
    setSelectedTodo(todo);
    setModalKey((prev) => prev + 1);
    setIsModalOpen(true);
  };

  const handleSave = (data: { title: string; description: string }) => {
    if (selectedTodo) {
      updateMutation.mutate({ id: selectedTodo.id, data });
    } else {
      createMutation.mutate(data);
    }
  };

  const handlePageChange = (page: number) => {
    router.push(`/?page=${page}`);
  };

  const handleDelete = (id: number) => {
    deleteMutation.mutate(id);
  };

  if (isLoading) {
    return (
      <div className="flex justify-center p-8">
        <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" />
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center p-8 text-destructive">
        Failed to load todos
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {data && data.data.length > 0 ? (
        data.data.map((todo) => (
          <TodoItem
            key={todo.id}
            todo={todo}
            onClick={handleOpenEdit}
            onDelete={handleDelete}
          />
        ))
      ) : (
        <div className="text-center text-muted-foreground p-8">
          No todos found. Add one!
        </div>
      )}

      {/* FAB */}
      <div className="fixed bottom-20 right-8 z-40">
        <Button
          size="icon"
          className="h-14 w-14 rounded-full shadow-lg transition-transform hover:scale-110"
          onClick={handleOpenCreate}
        >
          <Plus className="h-6 w-6" />
        </Button>
      </div>

      <TodoModal
        key={modalKey}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        todo={selectedTodo}
        onSave={handleSave}
        isSaving={createMutation.isPending || updateMutation.isPending}
      />

      {data && data.pagination.total_pages > 1 && (
        <PaginationControl
          totalPages={data.pagination.total_pages}
          currentPage={data.pagination.current_page}
          hasNext={data.pagination.has_next}
          hasPrev={data.pagination.has_prev}
          onPageChange={handlePageChange}
        />
      )}
    </div>
  );
}
