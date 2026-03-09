"use client";

import { api } from "@/lib/api";
import { showToast } from "@/lib/toast";
import {
  ListTodoFilterHistoriesResponse,
  Todo,
  TodoFilterHistoryQuery,
} from "@/types";
import { useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { z } from "zod";

const aiFilterSchema = z
  .string()
  .trim()
  .min(1, "入力してください")
  .max(50, "50文字以内で入力してください");

interface UseAIFilterProps {
  onFilterSuccess: (todos: Todo[]) => void;
  onClear: () => void;
}

export function useAIFilter({ onFilterSuccess, onClear }: UseAIFilterProps) {
  const [isFiltering, setIsFiltering] = useState(false);
  const [activeFilter, setActiveFilter] =
    useState<TodoFilterHistoryQuery | null>(null);
  const [filterHistories, setFilterHistories] = useState<
    TodoFilterHistoryQuery[]
  >([]);
  const router = useRouter();

  useEffect(() => {
    const fetchHistories = async () => {
      try {
        const res = await api.get<ListTodoFilterHistoriesResponse>(
          "/todo/filter_histories",
        );
        setFilterHistories(res.data.queries);
      } catch (error) {
        console.error("Failed to fetch filter histories", error);
      }
    };
    fetchHistories();
  }, []);

  const handleAIFilter = useCallback(
    async (query: string, currentPage: number) => {
      const validation = aiFilterSchema.safeParse(query);
      if (!validation.success) {
        showToast(validation.error.issues[0].message, "error");
        return;
      }

      setIsFiltering(true);
      try {
        const res = await api.get<Todo[]>("/todo/ai_filter", {
          params: { query: validation.data },
        });
        onFilterSuccess(res.data);
        setActiveFilter({ id: "custom", query: validation.data });
        if (currentPage !== 1) {
          router.push("/?page=1");
        }
      } catch (error) {
        console.error("AI filter failed", error);
        showToast("AIフィルターに失敗しました", "error");
      } finally {
        setIsFiltering(false);
      }
    },
    [onFilterSuccess, router],
  );

  const handleHistoryFilter = useCallback(
    async (history: TodoFilterHistoryQuery, currentPage: number) => {
      setIsFiltering(true);
      try {
        const res = await api.get<Todo[]>("/todo/filter_by_query_id", {
          params: { query_id: history.id },
        });
        onFilterSuccess(res.data);
        setActiveFilter(history);
        if (currentPage !== 1) {
          router.push("/?page=1");
        }
      } catch (error) {
        console.error("History filter failed", error);
        showToast("履歴からのフィルターに失敗しました", "error");
      } finally {
        setIsFiltering(false);
      }
    },
    [onFilterSuccess, router],
  );

  const clearFilter = useCallback(() => {
    setActiveFilter(null);
    onClear();
  }, [onClear]);

  return {
    isFiltering,
    activeFilter,
    filterHistories,
    handleAIFilter,
    handleHistoryFilter,
    clearFilter,
  };
}
