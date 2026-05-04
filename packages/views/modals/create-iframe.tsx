"use client";

import { useState, useRef } from "react";
import { X as XIcon } from "lucide-react";
import { useCreateIframe } from "@multica/core/iframes/mutations";
import { useWorkspacePaths } from "@multica/core/paths";
import { toast } from "sonner";
import { Dialog, DialogContent, DialogTitle } from "@multica/ui/components/ui/dialog";
import { Button } from "@multica/ui/components/ui/button";
import { EmojiPicker } from "@multica/ui/components/common/emoji-picker";
import { ContentEditor, type ContentEditorRef, TitleEditor } from "../editor";
import { Popover, PopoverTrigger, PopoverContent } from "@multica/ui/components/ui/popover";

export function CreateIframeModal({ onClose }: { onClose: () => void }) {
  const wsPaths = useWorkspacePaths();

  const [title, setTitle] = useState("");
  const descEditorRef = useRef<ContentEditorRef>(null);
  const [icon, setIcon] = useState<string | undefined>();
  const [iconPickerOpen, setIconPickerOpen] = useState(false);
  const [iframeUrl, setIframeUrl] = useState("");
  const [iframeScript, setIframeScript] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const createIframe = useCreateIframe();

  const handleSubmit = async () => {
    if (!title.trim() || submitting) return;
    setSubmitting(true);
    try {
      const data = {
        title: title.trim(),
        description: descEditorRef.current?.getMarkdown()?.trim() || undefined,
        icon,
        iframe_url: iframeUrl.trim() || undefined,
        iframe_script: iframeScript.trim() || undefined,
      };
      const iframe = await createIframe.mutateAsync(data);
      onClose();
      toast.success("IFRAME created");
      window.location.href = wsPaths.iframeDetail(iframe.id);
    } catch (error) {
      toast.error("Failed to create iframe");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog open onOpenChange={(v) => { if (!v) onClose(); }}>
      <DialogContent showCloseButton={false} className="p-0 gap-0 flex flex-col overflow-hidden max-w-2xl">
        <DialogTitle className="sr-only">New IFRAME</DialogTitle>

        <div className="flex items-center justify-between px-5 pt-3 pb-2 shrink-0">
          <div className="flex items-center gap-1.5 text-xs">
            <span className="font-medium">New IFRAME</span>
          </div>
          <button
            onClick={onClose}
            className="rounded-sm p-1.5 opacity-70 hover:opacity-100 hover:bg-accent/60 transition-all cursor-pointer"
          >
            <XIcon className="size-4" />
          </button>
        </div>

        <div className="px-5 pb-2 shrink-0">
          <Popover open={iconPickerOpen} onOpenChange={setIconPickerOpen}>
            <PopoverTrigger
              render={
                <button
                  type="button"
                  className="text-2xl cursor-pointer rounded-lg p-1 -ml-1 hover:bg-accent/60 transition-colors"
                  title="Choose icon"
                >
                  {icon || "🖼️"}
                </button>
              }
            />
            <PopoverContent align="start" className="w-auto p-0">
              <EmojiPicker
                onSelect={(emoji) => {
                  setIcon(emoji);
                  setIconPickerOpen(false);
                }}
              />
            </PopoverContent>
          </Popover>
          <TitleEditor
            autoFocus
            defaultValue=""
            placeholder="IFRAME title"
            className="text-lg font-semibold"
            onChange={(v) => setTitle(v)}
            onSubmit={handleSubmit}
          />
        </div>

        <div className="flex-1 min-h-0 overflow-y-auto px-5">
          <ContentEditor
            ref={descEditorRef}
            defaultValue=""
            placeholder="Add description..."
            debounceMs={500}
          />

          <div className="mt-4 space-y-4">
            <div>
              <label className="text-sm font-medium mb-1 block">Iframe URL</label>
              <input
                type="url"
                value={iframeUrl}
                onChange={(e) => setIframeUrl(e.target.value)}
                placeholder="https://example.com"
                className="w-full px-3 py-2 rounded-md border bg-background text-sm"
              />
            </div>

            <div>
              <label className="text-sm font-medium mb-1 block">Or Custom Script</label>
              <textarea
                value={iframeScript}
                onChange={(e) => setIframeScript(e.target.value)}
                placeholder="<script>...</script> or custom HTML"
                className="w-full px-3 py-2 rounded-md border bg-background text-sm min-h-[100px] font-mono"
              />
            </div>

            <p className="text-xs text-muted-foreground">
              Provide either a URL to embed or a custom script. If both are provided, the URL will be used.
            </p>
          </div>
        </div>

        <div className="flex items-center justify-end px-4 py-3 border-t shrink-0">
          <Button size="sm" onClick={handleSubmit} disabled={!title.trim() || submitting}>
            {submitting ? "Creating..." : "Create IFRAME"}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
