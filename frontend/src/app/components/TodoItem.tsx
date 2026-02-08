import { Card } from "@/app/components/ui/card";
import { Button } from "@/app/components/ui/button";
import { Trash2 } from "lucide-react";
import { Todo } from "@/types";

interface TodoItemProps {
  todo: Todo;
  onDelete: (id: number) => void;
  onClick: (todo: Todo) => void;
}

export function TodoItem({ todo, onDelete, onClick }: TodoItemProps) {
  return (
    <Card className="p-4 flex items-center justify-between hover:bg-muted/50 transition-colors">
      <div
        className="flex-1 cursor-pointer font-medium truncate mr-4"
        onClick={() => onClick(todo)}
      >
        {todo.title}
      </div>
      <Button
        variant="ghost"
        size="icon"
        onClick={(e) => {
          e.stopPropagation();
          onDelete(todo.id);
        }}
        className="text-muted-foreground hover:text-destructive shrink-0"
      >
        <Trash2 className="h-5 w-5" />
      </Button>
    </Card>
  );
}
