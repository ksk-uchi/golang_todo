import { api } from "@/lib/api";

export interface LoginRequest {
  email: string;
  password: string;
}

export const authService = {
  login: async (data: LoginRequest) => {
    const response = await api.post("/auth/login", data);
    return response.data;
  },
};
