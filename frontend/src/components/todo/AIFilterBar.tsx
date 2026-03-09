"use client";

import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";
import { TodoFilterHistoryQuery } from "@/types";
import { History, Search, Sparkles, X } from "lucide-react";
import { useState } from "react";

export interface AIFilterBarProps {
  filterHistories: TodoFilterHistoryQuery[];
  activeFilter: TodoFilterHistoryQuery | null;
  onSearchAIFilter: (query: string) => void;
  onSearchHistory: (history: TodoFilterHistoryQuery) => void;
  onClearFilter: () => void;
}

export function AIFilterBar({
  filterHistories,
  activeFilter,
  onSearchAIFilter,
  onSearchHistory,
  onClearFilter,
}: AIFilterBarProps) {
  const [inputValue, setInputValue] = useState("");
  const [isPopoverOpen, setIsPopoverOpen] = useState(false);

  const handleSearch = () => {
    onSearchAIFilter(inputValue);
    setIsPopoverOpen(false);
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleSearch();
    }
  };

  return (
    <div className="space-y-4">
      <div className="relative group">
        <div className="absolute -inset-0.5 bg-gradient-to-r from-purple-500 to-blue-500 rounded-lg blur opacity-25 group-hover:opacity-50 transition duration-1000 group-hover:duration-200"></div>
        <div className="relative bg-background rounded-lg">
          <Popover open={isPopoverOpen} onOpenChange={setIsPopoverOpen}>
            <PopoverTrigger asChild>
              <div className="relative">
                <Sparkles className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-blue-500" />
                <Input
                  className="pl-10 pr-10 border-blue-200 focus-visible:ring-blue-500"
                  placeholder="AIに頼んでタスクを絞り込む... (例: 「直近一週間以内に完了になったもの」)"
                  value={inputValue}
                  onChange={(e) => setInputValue(e.target.value)}
                  onKeyDown={handleKeyDown}
                  maxLength={50}
                />
                <button
                  onClick={handleSearch}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-blue-500 hover:text-blue-700 transition-colors"
                >
                  <Search className="w-4 h-4" />
                </button>
              </div>
            </PopoverTrigger>
            {filterHistories.length > 0 && (
              <PopoverContent
                className="w-[var(--radix-popover-trigger-width)] p-1"
                align="start"
              >
                <div className="text-xs font-medium text-muted-foreground px-2 py-1.5">
                  最近の検索
                </div>
                <div className="max-h-60 overflow-y-auto">
                  {filterHistories.map((history) => (
                    <div
                      key={history.id}
                      className="cursor-pointer hover:bg-muted p-2 text-sm rounded-md flex items-center gap-2"
                      onClick={() => {
                        onSearchHistory(history);
                        setIsPopoverOpen(false);
                        setInputValue(history.query);
                      }}
                    >
                      <History className="w-4 h-4 text-muted-foreground" />
                      <span className="truncate">{history.query}</span>
                    </div>
                  ))}
                </div>
              </PopoverContent>
            )}
          </Popover>
        </div>
      </div>

      {activeFilter && (
        <Badge
          variant="secondary"
          className="px-3 py-1 flex w-fit items-center gap-2 animate-in fade-in slide-in-from-top-1"
        >
          <Sparkles className="w-4 h-4 text-blue-500" />
          <span className="text-sm">AI抽出: {activeFilter.query}</span>
          <X
            className="w-4 h-4 cursor-pointer hover:text-destructive transition-colors ml-1"
            onClick={() => {
              onClearFilter();
              setInputValue("");
            }}
          />
        </Badge>
      )}
    </div>
  );
}
