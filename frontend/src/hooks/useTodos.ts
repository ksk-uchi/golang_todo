"use client";

import { api } from "@/lib/api";
import { showToast } from "@/lib/toast";
import { ListTodoResponse, Todo } from "@/types";
import { useRouter, useSearchParams } from "next/navigation";
import { useEffect, useState } from "react";

export function useTodos() {
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
  const [hideDone, setHideDone] = useState(true);

  useEffect(() => {
    const fetchTodos = async () => {
      setIsLoading(true);
      try {
        const res = await api.get<ListTodoResponse>("/todo", {
          params: { page: currentPage, include_done: !hideDone },
        });
        if (res.data.data.length === 0 && currentPage > 1) {
          router.push(`/?page=${res.data.pagination.total_pages}`);
          return;
        }
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
  }, [currentPage, router, hideDone]);

  const handleToggleDone = async (id: number, isDone: boolean) => {
    try {
      const res = await api.put(`/todo/${id}/done`, { is_done: isDone });
      const updatedTodo = res.data;
      setTodos((prev) =>
        prev.map((t) => (t.id === updatedTodo.id ? updatedTodo : t)),
      );
      showToast("Todo status updated", "success");
    } catch {
      showToast("Failed to update status", "error");
    }
  };

  const handleCreate = async (data: {
    title: string;
    description: string;
    isDone?: boolean;
  }) => {
    setIsSaving(true);
    try {
      const res = await api.post("/todo", {
        title: data.title,
        description: data.description,
      });
      let createdTodo = res.data;

      if (data.isDone) {
        const resDone = await api.put(`/todo/${createdTodo.id}/done`, {
          is_done: true,
        });
        createdTodo = resDone.data;
      }

      setTodos((prev) => [createdTodo, ...prev]);
      setIsModalOpen(false);
      showToast("Todo created successfully", "success");
    } catch {
      showToast("Failed to create todo", "error");
    } finally {
      setIsSaving(false);
    }
  };

  const handleUpdate = async ({
    id,
    data,
  }: {
    id: number;
    data: { title: string; description: string; isDone?: boolean };
  }) => {
    setIsSaving(true);
    try {
      let currentTodo = todos.find((t) => t.id === id);
      if (!currentTodo && selectedTodo?.id === id) currentTodo = selectedTodo!;
      if (!currentTodo) throw new Error("Todo not found");

      const isCurrentlyDone = !!currentTodo.done_at;
      const newIsDone =
        data.isDone !== undefined ? data.isDone : isCurrentlyDone;

      if (isCurrentlyDone && !newIsDone) {
        await api.put(`/todo/${id}/done`, { is_done: false });
      }

      const resPatch = await api.patch(`/todo/${id}`, {
        title: data.title,
        description: data.description,
      });
      let updatedTodo = resPatch.data;

      if (!isCurrentlyDone && newIsDone) {
        const resDone = await api.put(`/todo/${id}/done`, { is_done: true });
        updatedTodo = resDone.data;
      }

      setTodos((prev) =>
        prev.map((t) => (t.id === updatedTodo.id ? updatedTodo : t)),
      );
      setSelectedTodo(updatedTodo);
      showToast("Todo updated successfully", "success");
    } catch {
      showToast("Failed to update todo", "error");
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
      showToast("Todo deleted successfully", "success");
    } catch {
      showToast("Failed to delete todo", "error");
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

  const handleSave = (data: {
    title: string;
    description: string;
    isDone?: boolean;
  }) => {
    if (selectedTodo) {
      handleUpdate({ id: selectedTodo.id, data });
    } else {
      handleCreate(data);
    }
  };

  const handlePageChange = (page: number) => {
    router.push(`/?page=${page}`);
  };

  return {
    todos,
    pagination,
    isLoading,
    error,
    isSaving,
    isModalOpen,
    selectedTodo,
    modalKey,
    hideDone,
    setHideDone,
    handleToggleDone,
    handleDelete,
    handleOpenCreate,
    handleOpenEdit,
    handleSave,
    handlePageChange,
    setIsModalOpen,
  };
}
