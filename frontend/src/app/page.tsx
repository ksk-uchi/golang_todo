"use client";

import { AIFilterBar } from "@/components/todo/AIFilterBar";
import { PaginationControl } from "@/components/todo/PaginationControl";
import { TodoItem } from "@/components/todo/TodoItem";
import { TodoModal } from "@/components/todo/TodoModal";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Label } from "@/components/ui/label";
import { useAIFilter } from "@/hooks/useAIFilter";
import { useTodos } from "@/hooks/useTodos";
import { Plus } from "lucide-react";
import { useRouter, useSearchParams } from "next/navigation";
import { Suspense } from "react";

function HomeContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const pageParam = searchParams.get("page");
  const currentPage = pageParam ? parseInt(pageParam, 10) : 1;

  const {
    isFiltering,
    activeFilter,
    filterHistories,
    handleAIFilter,
    handleHistoryFilter,
    clearFilter,
  } = useAIFilter({
    onFilterSuccess: (todos) => setTodos(todos),
    onClear: () => fetchTodos(),
  });

  const {
    todos,
    setTodos,
    fetchTodos,
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
  } = useTodos(activeFilter !== null);

  if (isLoading || isFiltering) {
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
      <AIFilterBar
        activeFilter={activeFilter}
        filterHistories={filterHistories}
        onSearchAIFilter={(query) => handleAIFilter(query, currentPage)}
        onSearchHistory={(history) => handleHistoryFilter(history, currentPage)}
        onClearFilter={clearFilter}
      />

      {!activeFilter && (
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
      )}

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
          {activeFilter
            ? "該当するToDoは見つかりませんでした。"
            : "No todos found. Add one!"}
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

      {!activeFilter && pagination && pagination.total_pages > 1 && (
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
