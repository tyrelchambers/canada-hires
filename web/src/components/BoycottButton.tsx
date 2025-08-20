import { Button } from "@/components/ui/button";
import { useToggleBoycott, useBoycottStats } from "@/hooks/useBoycotts";
import { useCurrentUser } from "@/hooks/useAuth";
import { useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBan, faQuestion } from "@fortawesome/free-solid-svg-icons";
import { toast } from "sonner";

interface BoycottButtonProps {
  businessName: string;
  businessAddress: string;
  className?: string;
}

export function BoycottButton({
  businessName,
  businessAddress,
  className,
}: BoycottButtonProps) {
  const { data: user } = useCurrentUser();
  const [isLoading, setIsLoading] = useState(false);

  const { data: boycottStats, isLoading: statsLoading } = useBoycottStats(
    businessName,
    businessAddress,
  );

  const toggleMutation = useToggleBoycott();

  const handleToggle = async () => {
    if (!user) {
      toast.error("Please log in to boycott businesses");
      return;
    }

    setIsLoading(true);
    try {
      const result = await toggleMutation.mutateAsync({
        business_name: businessName,
        business_address: businessAddress,
      });

      if (result.is_boycotting) {
        toast.success("You are now boycotting this business");
      } else {
        toast.success("You are no longer boycotting this business");
      }
    } catch (error: any) {
      toast.error(error.response?.data?.message || "Failed to toggle boycott");
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className={`flex items-center gap-2 ${className}`}>
      <Button
        variant={boycottStats?.is_boycotted_by_user ? "secondary" : "outline"}
        onClick={handleToggle}
        disabled={isLoading || statsLoading || !user}
        className="w-full"
      >
        <FontAwesomeIcon
          icon={boycottStats?.is_boycotted_by_user ? faBan : faQuestion}
          className="mr-2"
        />
        {boycottStats?.is_boycotted_by_user
          ? "I'm boycotting!"
          : "Are you boycotting?"}
      </Button>
    </div>
  );
}
