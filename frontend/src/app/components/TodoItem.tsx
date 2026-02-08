import { Card, CardContent, CardAction } from "@/app/components/ui/card";
import { Button } from "@/app/components/ui/button";
import { Trash2 } from "lucide-react";
import { Todo } from "@/types";
import { cn } from "@/lib/utils";
import { TypographyP } from "@/app/components/ui/typography";

interface TodoItemProps {
  todo: Todo;
  onDelete: (id: number) => void;
  onClick: (todo: Todo) => void;
}

export function TodoItem({ todo, onDelete, onClick }: TodoItemProps) {
  return (
    <Card className="px-4 py-2 flex justify-between hover:bg-muted/50 transition-colors group min-h-[3rem]">
      <CardContent
        className={cn(
          "flex-1 cursor-pointer font-medium truncate mr-4 columns-2",
        )}
        onClick={() => onClick(todo)}
      >
        <TypographyP text={todo.title} />
        <CardAction className={cn("flex items-center")}>
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
