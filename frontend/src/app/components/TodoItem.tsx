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
    <Card className="px-4 py-2 flex items-center justify-between hover:bg-muted/50 transition-colors group min-h-[3rem]">
      <div
        className="flex-1 cursor-pointer font-medium truncate mr-4"
        onClick={() => onClick(todo)}
      >
        {todo.title}
      </div>
      <div className="opacity-0 group-hover:opacity-100 transition-opacity flex items-center">
        <Button
          variant="ghost"
          size="icon"
          onClick={(e) => {
            e.stopPropagation();
            onDelete(todo.id);
          }}
          className="text-muted-foreground hover:text-destructive h-8 w-8"
        >
          <Trash2 className="h-4 w-4" />
        </Button>
      </div>
    </Card>
  );
}
