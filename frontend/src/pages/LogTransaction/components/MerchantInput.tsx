import { Input } from "@/components/ui/Input";
import { Store } from "lucide-react";

export function MerchantInput({ value, onChange }: { value: string; onChange: (v: string) => void }) {
  return (
    <div className="relative">
      <Store size={14} className="absolute left-3 top-9 text-muted-foreground" />
      <Input label="Merchant" value={value} onChange={(e) => onChange(e.target.value)} className="pl-9" />
    </div>
  );
}
export default MerchantInput;
