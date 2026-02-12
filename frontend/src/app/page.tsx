"use client";

import { PaginationControl } from "@/app/components/PaginationControl";
import { TodoItem } from "@/app/components/TodoItem";
import { TodoModal } from "@/app/components/TodoModal";
import { Button } from "@/app/components/ui/button";
import { api } from "@/lib/api";
import { ListTodoResponse, Todo } from "@/types";
import { Plus } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function Home() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const pageParam = searchParams.get("page");
  const currentPage = pageParam ? parseInt(pageParam, 10) : 1;

  const [todos, setTodos] = useState<Todo[]>([]);
  const [pagination, setPagination] = useState<
    ListTodoResponse["pagination"] | null
  >(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);
  const [isSaving, setIsSaving] = useState(false);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null);
  const [modalKey, setModalKey] = useState(0);

  useEffect(() => {
    const fetchTodos = async () => {
      setIsLoading(true);
      try {
        const res = await api.get<ListTodoResponse>("/todo", {
          params: { page: currentPage },
        });
        setTodos(res.data.data);
        setPagination(res.data.pagination);
        setError(null);
      } catch (err) {
        setError(err as Error);
      } finally {
        setIsLoading(false);
      }
    };
    fetchTodos();
  }, [currentPage]);

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

  const handleCreate = async (data: { title: string; description: string }) => {
    setIsSaving(true);
    try {
      const res = await api.post("/todo", data);
      setTodos((prev) => [res.data, ...prev]);
      setIsModalOpen(false);
      toastHandler("Todo created successfully", "success");
    } catch {
      toastHandler("Failed to create todo", "error");
    } finally {
      setIsSaving(false);
    }
  };

  const handleUpdate = async ({
    id,
    data,
  }: {
    id: number;
    data: { title: string; description: string };
  }) => {
    setIsSaving(true);
    try {
      const res = await api.patch(`/todo/${id}`, data);
      const updatedTodo = res.data;
      setTodos((prev) =>
        prev.map((t) => (t.id === updatedTodo.id ? updatedTodo : t)),
      );
      setSelectedTodo(updatedTodo);
      toastHandler("Todo updated successfully", "success");
    } catch {
      toastHandler("Failed to update todo", "error");
    } finally {
      setIsSaving(false);
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await api.delete(`/todo/${id}`);
      if (todos.length === 1 && currentPage > 1) {
        router.push(`/?page=${currentPage - 1}`);
      } else {
        setTodos((prev) => prev.filter((t) => t.id !== id));
      }
      toastHandler("Todo deleted successfully", "success");
    } catch {
      toastHandler("Failed to delete todo", "error");
    }
  };

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
      handleUpdate({ id: selectedTodo.id, data });
    } else {
      handleCreate(data);
    }
  };

  const handlePageChange = (page: number) => {
    router.push(`/?page=${page}`);
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
      {todos.length > 0 ? (
        todos.map((todo) => (
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
        isSaving={isSaving}
      />

      {pagination && pagination.total_pages > 1 && (
        <PaginationControl
          totalPages={pagination.total_pages}
          currentPage={pagination.current_page}
          hasNext={pagination.has_next}
          hasPrev={pagination.has_prev}
          onPageChange={handlePageChange}
        />
      )}
    </div>
  );
}
