"use client";

import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Todo } from "@/types";
import { TodoItem } from "@/app/components/TodoItem";
import { TodoModal } from "@/app/components/TodoModal";
import { Button } from "@/app/components/ui/button";
import { Plus } from "lucide-react";
import { toast } from "sonner";
import { useState } from "react";

export default function Home() {
  const queryClient = useQueryClient();
  const {
    data: todos,
    isLoading,
    error,
  } = useQuery<Todo[]>({
    queryKey: ["todos"],
    queryFn: async () => {
      const res = await api.get("/todo");
      return res.data;
    },
  });

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedTodo, setSelectedTodo] = useState<Todo | null>(null);

  // Mutations
  const createMutation = useMutation({
    mutationFn: async (data: { title: string; description: string }) => {
      const res = await api.post("/todo", data);
      return res.data;
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      setIsModalOpen(false);
      toast.success("Todo created successfully", { duration: 6000 });
    },
    onError: () => {
      toast.error("Failed to create todo", { duration: 6000 });
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
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      // Update selected todo if currently viewing it to reflect changes (e.g. updated_at if displayed)
      setSelectedTodo(updatedTodo);
      // Modal stays open
      toast.success("Todo updated successfully", { duration: 6000 });
    },
    onError: () => {
      toast.error("Failed to update todo", { duration: 6000 });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: async (id: number) => {
      await api.delete(`/todo/${id}`);
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["todos"] });
      toast.success("Todo deleted successfully", { duration: 6000 });
    },
    onError: () => {
      toast.error("Failed to delete todo", { duration: 6000 });
    },
  });

  const handleOpenCreate = () => {
    setSelectedTodo(null);
    setIsModalOpen(true);
  };

  const handleOpenEdit = (todo: Todo) => {
    setSelectedTodo(todo);
    setIsModalOpen(true);
  };

  const handleSave = (data: { title: string; description: string }) => {
    if (selectedTodo) {
      updateMutation.mutate({ id: selectedTodo.id, data });
    } else {
      createMutation.mutate(data);
    }
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
    <div className="space-y-4 pb-20">
      {todos && todos.length > 0 ? (
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
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        todo={selectedTodo}
        onSave={handleSave}
        isSaving={createMutation.isPending || updateMutation.isPending}
      />
    </div>
  );
}
