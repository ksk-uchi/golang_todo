"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api";
import { Todo } from "@/types";
import { TodoItem } from "@/app/components/TodoItem";

export default function Home() {
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
            onClick={(t) => console.log("Clicked:", t)}
            onDelete={(id) => console.log("Delete:", id)}
          />
        ))
      ) : (
        <div className="text-center text-muted-foreground p-8">
          No todos found. Add one!
        </div>
      )}
    </div>
  );
}
