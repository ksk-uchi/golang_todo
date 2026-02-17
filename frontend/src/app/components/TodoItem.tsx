import { Button } from "@/app/components/ui/button";
import { Card, CardAction, CardContent } from "@/app/components/ui/card";
import { Checkbox } from "@/app/components/ui/checkbox";
import { TypographyP } from "@/app/components/ui/typography";
import { cn } from "@/lib/utils";
import { Todo } from "@/types";
import { Trash2 } from "lucide-react";

interface TodoItemProps {
  todo: Todo;
  onDelete: (id: number) => void;
  onClick: (todo: Todo) => void;
  onToggleDone: (id: number, isDone: boolean) => void;
}

export function TodoItem({
  todo,
  onDelete,
  onClick,
  onToggleDone,
}: TodoItemProps) {
  const isDone = !!todo.done_at;

  return (
    <Card className="px-4 py-2 flex justify-between hover:bg-muted/50 transition-colors group min-h-[3rem]">
      <CardContent className="flex items-center flex-1 mr-4 overflow-hidden">
        <Checkbox
          checked={isDone}
          onCheckedChange={(checked) => onToggleDone(todo.id, checked === true)}
          className="mr-4"
          onClick={(e) => e.stopPropagation()}
        />
        <div
          className={cn(
            "flex-1 cursor-pointer font-medium truncate",
            isDone && "text-muted-foreground line-through",
          )}
          onClick={() => onClick(todo)}
        >
          <TypographyP text={todo.title} />
        </div>
        <CardAction className={cn("flex items-center ml-auto")}>
          <Button
            variant="ghost"
            size="icon"
            className={cn(
              "text-muted-foreground hover:text-destructive cursor-pointer h-8 w-8",
            )}
            onClick={(e) => {
              e.stopPropagation();
              onDelete(todo.id);
            }}
          >
            <Trash2 className="h-4 w-4" />
          </Button>
        </CardAction>
      </CardContent>
    </Card>
  );
}
