import { AstroProviders } from "@/components/astro-providers";

export default function LandingPage() {
  return (
    <AstroProviders>
      <div className="flex items-center justify-center h-screen">
        <div className="text-center">
          <h1 className="text-4xl font-bold mb-4">Multica</h1>
          <p className="text-muted-foreground">Project Management for Human + Agent Teams</p>
        </div>
      </div>
    </AstroProviders>
  );
}
