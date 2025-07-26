import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faChevronLeft, faChevronRight } from "@fortawesome/free-solid-svg-icons";
import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  totalItems: number;
  itemsPerPage: number;
  onPageChange: (page: number) => void;
  className?: string;
}

export function Pagination({
  currentPage,
  totalPages,
  totalItems,
  itemsPerPage,
  onPageChange,
  className,
}: PaginationProps) {
  const startItem = (currentPage - 1) * itemsPerPage + 1;
  const endItem = Math.min(currentPage * itemsPerPage, totalItems);

  // Calculate page numbers to show - responsive based on screen size
  const getVisiblePages = (isMobile = false) => {
    const delta = isMobile ? 1 : 2; // Fewer pages on mobile
    const maxPages = isMobile ? 5 : 7; // Fewer total pages on mobile
    const pages: (number | string)[] = [];
    
    if (totalPages <= maxPages) {
      // Show all pages if within limit
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      // Always show first page
      pages.push(1);
      
      if (currentPage > delta + 2) {
        pages.push("...");
      }
      
      // Show pages around current page
      const start = Math.max(2, currentPage - delta);
      const end = Math.min(totalPages - 1, currentPage + delta);
      
      for (let i = start; i <= end; i++) {
        pages.push(i);
      }
      
      if (currentPage < totalPages - delta - 1) {
        pages.push("...");
      }
      
      // Always show last page
      if (totalPages > 1) {
        pages.push(totalPages);
      }
    }
    
    return pages;
  };

  const desktopPages = getVisiblePages(false);
  const mobilePages = getVisiblePages(true);

  return (
    <div className={cn("space-y-4", className)}>
      {/* Results info - always show */}
      <div className="text-sm text-muted-foreground text-center sm:text-left">
        Showing {startItem} to {endItem} of {totalItems} results
      </div>
      
      {/* Desktop pagination */}
      <div className="hidden sm:flex items-center justify-center space-x-2">
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(currentPage - 1)}
          disabled={currentPage <= 1}
        >
          <FontAwesomeIcon icon={faChevronLeft} className="h-4 w-4" />
          Previous
        </Button>
        
        {desktopPages.map((page, index) => (
          <Button
            key={index}
            variant={page === currentPage ? "default" : "outline"}
            size="sm"
            onClick={() => typeof page === "number" && onPageChange(page)}
            disabled={typeof page !== "number"}
            className={cn(
              "min-w-[40px]",
              typeof page !== "number" && "cursor-default"
            )}
          >
            {page}
          </Button>
        ))}
        
        <Button
          variant="outline"
          size="sm"
          onClick={() => onPageChange(currentPage + 1)}
          disabled={currentPage >= totalPages}
        >
          Next
          <FontAwesomeIcon icon={faChevronRight} className="h-4 w-4" />
        </Button>
      </div>

      {/* Mobile pagination */}
      <div className="sm:hidden space-y-3">
        {/* Previous/Next buttons */}
        <div className="flex justify-between">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage <= 1}
            className="flex-1 mr-2"
          >
            <FontAwesomeIcon icon={faChevronLeft} className="h-4 w-4 mr-1" />
            Previous
          </Button>
          
          <Button
            variant="outline"
            size="sm"
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage >= totalPages}
            className="flex-1 ml-2"
          >
            Next
            <FontAwesomeIcon icon={faChevronRight} className="h-4 w-4 ml-1" />
          </Button>
        </div>
        
        {/* Page numbers - compact layout */}
        <div className="flex items-center justify-center space-x-1">
          {mobilePages.map((page, index) => (
            <Button
              key={index}
              variant={page === currentPage ? "default" : "outline"}
              size="sm"
              onClick={() => typeof page === "number" && onPageChange(page)}
              disabled={typeof page !== "number"}
              className={cn(
                "min-w-[32px] h-8 px-2 text-xs",
                typeof page !== "number" && "cursor-default"
              )}
            >
              {page}
            </Button>
          ))}
        </div>
        
        {/* Current page indicator */}
        <div className="text-xs text-muted-foreground text-center">
          Page {currentPage} of {totalPages}
        </div>
      </div>
    </div>
  );
}