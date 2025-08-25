import { Button } from "../ui/button";
import { Link } from "@tanstack/react-router";

import img1 from "@/assets/heatmap.png";
import img2 from "@/assets/non comp.png";
import img3 from "@/assets/reports.png";
import { Squares } from "../ui/squares-background";
import { Input } from "../ui/input";
import { useState } from "react";
import { useSendLoginLink } from "@/hooks/useAuth";
import clsx from "clsx";
import { toast } from "sonner";

export default function Hero() {
  const [email, setEmail] = useState("");
  const sendLoginLinkMutation = useSendLoginLink();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    sendLoginLinkMutation.mutate(email, {
      onSuccess: () => {
        setEmail("");
        toast.success("Login link sent successfully!");
      },
    });
  };

  return (
    <section className="bg-white border-b border-border relative  h-full lg:h-[800px] flex justify-center">
      <div className="max-w-screen-2xl mx-auto w-full border-x  border-border z-10 relative p-4 lg:p-20 h-full flex gap-6 lg:gap-20 items-center ">
        <div className="flex flex-col">
          <h1 className="text-3xl lg:text-6xl -tracking-[0.015em] font-bold mb-6 text-raising-black">
            Selling Out Canadian Jobs
          </h1>
          <p className="md:text-xl text-gray-500 font-light  max-w-3xl">
            The Temporary Foreign Worker (TFW) program is meant to fill labour
            shortages, but some companies exploit it to hire cheaper foreign
            labour instead of Canadians. We track the data so you can see which
            companies are abusing the system and choose where you spend your
            money.
          </p>

          <div className="bg-white  border border-border my-6 rounded-md overflow-hidden">
            <div className="p-4 md:p-8">
              <p className="mb-2 text-foreground font-medium">
                Sign in to add your reports
              </p>

              <form onSubmit={handleSubmit}>
                <div className="flex gap-4 flex-col md:flex-row">
                  <Input
                    placeholder="cold@canofbeans.com"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                    disabled={sendLoginLinkMutation.isPending}
                  />{" "}
                  <Button
                    disabled={sendLoginLinkMutation.isPending}
                    className={clsx(
                      sendLoginLinkMutation.isPending && "opacity-50",
                    )}
                  >
                    {sendLoginLinkMutation.isPending
                      ? "Sending..."
                      : "Get sign in link"}
                  </Button>
                </div>
              </form>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 h-32 md:h-16 border-border border-t">
              <Link
                to="/lmia"
                className="flex items-center justify-center bg-white hover:bg-secondary transition-all  "
              >
                Search LMIA records
              </Link>
              <Link
                to="/jobs"
                className="flex items-center justify-center bg-space-cadet text-isabelline hover:bg-space-cadet/90 transition-all"
              >
                Browse Job Postings
              </Link>
            </div>
          </div>
        </div>
        <div className="grid-cols-1 w-[800px] hidden lg:grid ">
          <div className="grid grid-cols-2 w-full ">
            <img src={img1} className="w-full" />

            <img src={img2} className="w-full" />
          </div>
          <img src={img3} className="w-full" />
        </div>
      </div>

      <Squares
        direction="down"
        speed={0.3}
        squareSize={100}
        borderColor="#CCC"
        hoverFillColor="#222"
        className="!bg-white absolute inset-0"
      />
    </section>
  );
}
