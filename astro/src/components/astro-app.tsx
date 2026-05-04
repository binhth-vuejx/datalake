import { AstroProviders } from "./astro-providers";

export default function AstroApp() {
  return (
    <AstroProviders>
      <div className="space-y-8">
        <section>
          <h2 className="text-2xl font-bold mb-4">React Integration Demo</h2>
          <p className="text-muted-foreground">
            This React component is wrapped in AstroProviders and demonstrates
            successful React integration with Astro, Vue, and Svelte.
          </p>
        </section>

        <section>
          <h3 className="text-lg font-semibold mb-2">Features</h3>
          <ul className="list-disc list-inside space-y-1 text-muted-foreground">
            <li>React components from @multica/views</li>
            <li>Shared CoreProvider for API and auth</li>
            <li>Navigation via History API</li>
            <li>WebSocket support</li>
          </ul>
        </section>
      </div>
    </AstroProviders>
  );
}
