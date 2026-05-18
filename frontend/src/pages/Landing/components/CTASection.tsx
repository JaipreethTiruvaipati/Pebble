import { Link } from "@tanstack/react-router";
import { Button } from "@/components/ui/Button";
import { ROUTES } from "@/routes";
import { ArrowRight } from "lucide-react";

export function CTASection() {
  return (
    <section className="mx-auto max-w-6xl px-6 pb-32">
      <div className="relative overflow-hidden rounded-3xl border border-border bg-card px-10 py-20 text-center">
        <div className="absolute inset-0 mesh-gradient opacity-60" />
        <div className="relative">
          <h2 className="font-display text-4xl md:text-5xl font-semibold tracking-tight max-w-2xl mx-auto">
            Stop leaking money. Start compounding it.
          </h2>
          <p className="mx-auto mt-4 max-w-xl text-muted-foreground">
            Join 12,000+ Indians turning every impulse into a long-term position.
          </p>
          <div className="mt-8 flex flex-wrap items-center justify-center gap-3">
            <Button asChild size="lg" variant="coral"><Link to={ROUTES.signup}>Start free <ArrowRight size={16} /></Link></Button>
            <Button asChild size="lg" variant="outline"><Link to={ROUTES.dashboard}>Tour the app</Link></Button>
          </div>
        </div>
      </div>
    </section>
  );
}
export default CTASection;

