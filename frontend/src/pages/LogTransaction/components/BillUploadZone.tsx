import { useDropzone } from "react-dropzone";
import { motion } from "framer-motion";
import { UploadCloud, FileImage } from "lucide-react";

export function BillUploadZone({ onFile }: { onFile: (file: File) => void }) {
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    accept: { "image/*": [], "application/pdf": [] },
    multiple: false,
    onDrop: (files) => files[0] && onFile(files[0]),
  });

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }} animate={{ opacity: 1, y: 0 }}
      {...(getRootProps() as any)}
      className={`relative cursor-pointer overflow-hidden rounded-3xl border-2 border-dashed p-16 text-center transition-all ${
        isDragActive ? "border-coral bg-coral/10 scale-[1.01]" : "border-border bg-card/40 hover:border-coral/60 hover:bg-coral/5"
      }`}
    >
      <input {...getInputProps()} />
      <div className="mx-auto grid h-20 w-20 place-items-center rounded-2xl bg-card-elevated text-coral">
        <UploadCloud size={32} />
      </div>
      <h3 className="mt-6 text-2xl">Drop your bill here</h3>
      <p className="mt-2 text-sm text-muted-foreground">PDF, PNG, JPG · max 10MB</p>
      <button className="mt-6 inline-flex items-center gap-2 rounded-xl bg-coral px-5 py-3 text-sm font-medium shadow-lg shadow-coral/20">
        <FileImage size={14} /> Browse files
      </button>
      <p className="mt-6 text-[11px] text-muted-foreground">
        Tip: forward Amazon / Zomato emails to <span className="font-mono text-foreground">bills@pebble.in</span>
      </p>
    </motion.div>
  );
}
export default BillUploadZone;
