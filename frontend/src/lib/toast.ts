import { toast } from "sonner";

export const showToast = (message: string, type: "success" | "error") => {
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
