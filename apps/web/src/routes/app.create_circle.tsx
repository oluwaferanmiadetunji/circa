import { createFileRoute } from "@tanstack/react-router";

import CreateGroup from "@/lib/pages/create-circle";

export const Route = createFileRoute("/app/create_circle")({
  component: CreateGroup,
});
