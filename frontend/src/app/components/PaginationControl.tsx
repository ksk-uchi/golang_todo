import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/app/components/ui/pagination";

interface PaginationControlProps {
  totalPages: number;
  currentPage: number;
  hasNext: boolean;
  hasPrev: boolean;
  onPageChange: (page: number) => void;
}

export function PaginationControl({
  totalPages,
  currentPage,
  hasNext,
  hasPrev,
  onPageChange,
}: PaginationControlProps) {
  const getPageNumbers = () => {
    const pages: (number | "ellipsis")[] = [];
    const maxDisplayed = 3;

    let startPage = Math.max(1, currentPage - 1);
    let endPage = Math.min(totalPages, currentPage + 1);

    if (currentPage === 1) {
      endPage = Math.min(totalPages, startPage + maxDisplayed - 1);
    } else if (currentPage === totalPages) {
      startPage = Math.max(1, endPage - maxDisplayed + 1);
    }

    // Ensure we always have at most 3 pages displayed if possible
    if (endPage - startPage + 1 < maxDisplayed) {
      if (startPage === 1) {
        endPage = Math.min(totalPages, startPage + maxDisplayed - 1);
      } else if (endPage === totalPages) {
        startPage = Math.max(1, endPage - maxDisplayed + 1);
      }
    }

    if (startPage > 1) {
      pages.push("ellipsis");
    }

    for (let i = startPage; i <= endPage; i++) {
      pages.push(i);
    }

    if (endPage < totalPages) {
      pages.push("ellipsis");
    }

    return pages;
  };

  const pages = getPageNumbers();

  const handlePageChange = (page: number) => {
    if (page !== currentPage) {
      onPageChange(page);
    }
  };

  return (
    <Pagination className="mt-4">
      <PaginationContent>
        <PaginationItem>
          <PaginationPrevious
            href="#"
            onClick={(e) => {
              e.preventDefault();
              if (hasPrev) handlePageChange(currentPage - 1);
            }}
            className={!hasPrev ? "pointer-events-none opacity-50" : ""}
            aria-disabled={!hasPrev}
          />
        </PaginationItem>

        {pages.map((page, index) => (
          <PaginationItem key={index}>
            {page === "ellipsis" ? (
              <PaginationEllipsis />
            ) : (
              <PaginationLink
                href="#"
                isActive={page === currentPage}
                onClick={(e) => {
                  e.preventDefault();
                  handlePageChange(page);
                }}
              >
                {page}
              </PaginationLink>
            )}
          </PaginationItem>
        ))}

        <PaginationItem>
          <PaginationNext
            href="#"
            onClick={(e) => {
              e.preventDefault();
              if (hasNext) handlePageChange(currentPage + 1);
            }}
            className={!hasNext ? "pointer-events-none opacity-50" : ""}
            aria-disabled={!hasNext}
          />
        </PaginationItem>
      </PaginationContent>
    </Pagination>
  );
}
