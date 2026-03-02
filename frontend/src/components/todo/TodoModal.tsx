"use client";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Todo } from "@/types";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { Controller, useForm, useWatch } from "react-hook-form";
import { z } from "zod";

const todoSchema = z.object({
  title: z.string().min(1, "Title is required").max(100, "Title is too long"),
  description: z.string().max(1000, "Description is too long"),
  isDone: z.boolean(),
});

type TodoFormValues = z.infer<typeof todoSchema>;

interface TodoModalProps {
  isOpen: boolean;
  onClose: () => void;
  todo?: Todo | null;
  onSave: (data: TodoFormValues) => void;
  isSaving?: boolean;
}

export function TodoModal({
  isOpen,
  onClose,
  todo,
  onSave,
  isSaving,
}: TodoModalProps) {
  const isEdit = !!todo;

  const {
    register,
    handleSubmit,
    control,
    reset,
    formState: { errors },
  } = useForm<TodoFormValues>({
    resolver: zodResolver(todoSchema),
    defaultValues: {
      title: todo?.title || "",
      description: todo?.description || "",
      isDone: !!todo?.done_at,
    },
  });

  useEffect(() => {
    if (isOpen) {
      reset({
        title: todo?.title || "",
        description: todo?.description || "",
        isDone: !!todo?.done_at,
      });
    }
  }, [isOpen, todo, reset]);

  const onSubmit = (data: TodoFormValues) => {
    onSave(data);
  };

  const isDoneValue = useWatch({ control, name: "isDone" });

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>{isEdit ? "Edit ToDo" : "Create ToDo"}</DialogTitle>
        </DialogHeader>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 py-4">
          <div className="flex items-center space-x-2">
            <Controller
              name="isDone"
              control={control}
              render={({ field }) => (
                <Checkbox
                  id="modal-done"
                  checked={field.value}
                  onCheckedChange={field.onChange}
                />
              )}
            />
            <Label htmlFor="modal-done">Done</Label>
          </div>
          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input
              id="title"
              placeholder="Title"
              {...register("title")}
              disabled={isDoneValue}
            />
            {errors.title && (
              <p className="text-sm text-destructive">{errors.title.message}</p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="description">Description</Label>
            <Textarea
              id="description"
              placeholder="Description"
              className="min-h-[100px]"
              {...register("description")}
              disabled={isDoneValue}
            />
            {errors.description && (
              <p className="text-sm text-destructive">
                {errors.description.message}
              </p>
            )}
          </div>
          <DialogFooter>
            <Button type="submit" disabled={isSaving}>
              {isSaving ? "Saving..." : isEdit ? "Save" : "Create"}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}
