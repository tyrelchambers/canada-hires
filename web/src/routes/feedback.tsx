import { createFileRoute } from "@tanstack/react-router";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Mail } from "lucide-react";
import { AuthNav } from "@/components/AuthNav";

export const Route = createFileRoute("/feedback")({
  component: FeedbackPage,
});

function FeedbackPage() {
  return (
    <section>
      <AuthNav />
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-2xl mx-auto">
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl">
                Feedback & Data Accuracy
              </CardTitle>
              <CardDescription>
                Help us improve JobWatch Canada by reporting data issues or sharing
                suggestions
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <div className="space-y-4">
                <h3 className="text-lg font-semibold">
                  We'd love to hear from you if you notice:
                </h3>
                <ul className="list-disc list-inside space-y-2 text-gray-700">
                  <li>Incorrect business information or data</li>
                  <li>Missing businesses or outdated listings</li>
                  <li>Technical issues with the platform</li>
                  <li>Suggestions for new features or improvements</li>
                  <li>Questions about our data sources or methodology</li>
                </ul>
              </div>

              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="font-semibold mb-3">Contact Us</h3>
                <p className="text-gray-700 mb-4">
                  Please send your feedback, questions, or data corrections to:
                </p>
                <Button asChild className="w-full sm:w-auto">
                  <a
                    href="mailto:connect@jobwatchcanada.com"
                    className="flex items-center gap-2"
                  >
                    <Mail className="h-4 w-4" />
                    connect@jobwatchcanada.com
                  </a>
                </Button>
              </div>

              <div className="text-sm text-gray-600">
                <p>
                  We appreciate your feedback and will work promptly to address
                  any issues. Your input helps us maintain the accuracy and
                  usefulness of this platform for all Canadians.
                </p>
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </section>
  );
}
