import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/app/components/ui/dialog";
import { Input } from "@/app/components/ui/input";
import { Textarea } from "@/app/components/ui/textarea";
import { Button } from "@/app/components/ui/button";
import { Todo } from "@/types";

interface TodoModalProps {
  isOpen: boolean;
  onClose: () => void;
  todo?: Todo | null; // if null, create mode
  onSave: (data: { title: string; description: string }) => void;
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

  const handleSubmit = () => {
    onSave({ title, description });
  };

  const isEdit = !!todo;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit ToDo" : "Create ToDo"}</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <Input
            id="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Title"
            className="col-span-3"
          />
          <Textarea
            id="description"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            placeholder="Description"
            className="col-span-3 min-h-[100px]"
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
