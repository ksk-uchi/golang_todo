import { Button } from "@/app/components/ui/button";
import { Checkbox } from "@/app/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/app/components/ui/dialog";
import { Input } from "@/app/components/ui/input";
import { Label } from "@/app/components/ui/label";
import { Textarea } from "@/app/components/ui/textarea";
import { Todo } from "@/types";
import { useState } from "react";

interface TodoModalProps {
  isOpen: boolean;
  onClose: () => void;
  todo?: Todo | null; // if null, create mode
  onSave: (data: {
    title: string;
    description: string;
    isDone?: boolean;
  }) => void;
  isSaving?: boolean;
}

export function TodoModal({
  isOpen,
  onClose,
  todo,
  onSave,
  isSaving,
}: TodoModalProps) {
  // Initialize state from props. Parent should define `key` to reset state when todo changes.
  const [title, setTitle] = useState(todo?.title || "");
  const [description, setDescription] = useState(todo?.description || "");
  const [isDone, setIsDone] = useState(!!todo?.done_at);

  const handleSubmit = () => {
    onSave({ title, description, isDone });
  };

  const isEdit = !!todo;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit ToDo" : "Create ToDo"}</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="flex items-center space-x-2">
            <Checkbox
              id="modal-done"
              checked={isDone}
              onCheckedChange={(c) => setIsDone(c === true)}
            />
            <Label htmlFor="modal-done">Done</Label>
          </div>
          <Input
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Title"
            className="col-span-3"
            disabled={isDone}
          />
          <Textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Description"
            className="col-span-3 min-h-[100px]"
            disabled={isDone}
          />
        </div>
        <DialogFooter>
          <Button onClick={handleSubmit} disabled={isSaving}>
            {isSaving ? "Saving..." : isEdit ? "Save" : "Create"}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
