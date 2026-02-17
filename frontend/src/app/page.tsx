"use client";

import { PaginationControl } from "@/app/components/PaginationControl";
import { TodoItem } from "@/app/components/TodoItem";
import { TodoModal } from "@/app/components/TodoModal";
import { Button } from "@/app/components/ui/button";
import { Checkbox } from "@/app/components/ui/checkbox";
import { Label } from "@/app/components/ui/label";
import { api } from "@/lib/api";
import { ListTodoResponse, Todo } from "@/types";
import { Plus } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense, useEffect, useState } from "react";
import { toast } from "sonner";

function HomeContent() {
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

  const handleToggleDone = async (id: number, isDone: boolean) => {
    try {
      const res = await api.put(`/todo/${id}/done`, { is_done: isDone });
      const updatedTodo = res.data;
      setTodos((prev) =>
        prev.map((t) => (t.id === updatedTodo.id ? updatedTodo : t)),
      );
      toastHandler("Todo status updated", "success");
    } catch {
      toastHandler("Failed to update status", "error");
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

      // 1. If becoming active (Done -> Active), toggle first
      if (isCurrentlyDone && !newIsDone) {
        await api.put(`/todo/${id}/done`, { is_done: false });
      }

      // 2. Update content
      // Note: If todo is Done and we are NOT changing it to Active, the backend would block PATCH.
      // However, the UI disables inputs when Done, ensuring title/desc shouldn't change generally.
      // But if we perform step 1, it is now Active, so PATCH works.
      const resPatch = await api.patch(`/todo/${id}`, {
        title: data.title,
        description: data.description,
      });
      let updatedTodo = resPatch.data;

      // 3. If becoming done (Active -> Done), toggle last
      if (!isCurrentlyDone && newIsDone) {
        const resDone = await api.put(`/todo/${id}/done`, { is_done: true });
        updatedTodo = resDone.data;
      }

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
      <div className="flex items-center space-x-2 pb-2">
        <Checkbox
          id="hide-done"
          checked={hideDone}
          onCheckedChange={(checked) => {
            const isChecked = checked === true;
            setHideDone(isChecked);
            router.push("/?page=1");
          }}
        />
        <Label
          htmlFor="hide-done"
          className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
        >
          完了したものを非表示にする
        </Label>
      </div>

      {todos.length > 0 ? (
        todos.map((todo) => (
          <TodoItem
            key={todo.id}
            todo={todo}
            onClick={handleOpenEdit}
            onDelete={handleDelete}
            onToggleDone={handleToggleDone}
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

export default function Home() {
  return (
    <Suspense
      fallback={
        <div className="flex justify-center p-8">
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full" />
        </div>
      }
    >
      <HomeContent />
    </Suspense>
  );
}
