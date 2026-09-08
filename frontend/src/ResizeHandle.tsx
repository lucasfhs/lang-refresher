import { useRef, useState, type KeyboardEvent, type PointerEvent } from "react";

interface Props {
  orientation: "vertical" | "horizontal";
  ariaLabel: string;
  onDrag: (delta: number) => void;
  onReset: () => void;
}

export function ResizeHandle({ orientation, ariaLabel, onDrag, onReset }: Props) {
  const lastPosition = useRef<number | null>(null);
  const [dragging, setDragging] = useState(false);

  const coordinate = (event: PointerEvent<HTMLDivElement>) => orientation === "vertical" ? event.clientX : event.clientY;

  const finishDragging = (event: PointerEvent<HTMLDivElement>) => {
    if (lastPosition.current === null) return;
    lastPosition.current = null;
    setDragging(false);
    document.body.classList.remove("is-resizing", `resize-${orientation}`);
    if (event.currentTarget.hasPointerCapture(event.pointerId)) {
      event.currentTarget.releasePointerCapture(event.pointerId);
    }
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLDivElement>) => {
    const delta = orientation === "vertical"
      ? event.key === "ArrowLeft" ? -16 : event.key === "ArrowRight" ? 16 : 0
      : event.key === "ArrowUp" ? -16 : event.key === "ArrowDown" ? 16 : 0;
    if (delta !== 0) {
      event.preventDefault();
      onDrag(delta);
    }
  };

  return (
    <div
      className={`resize-handle resize-handle-${orientation}${dragging ? " dragging" : ""}`}
      role="separator"
      aria-label={ariaLabel}
      aria-orientation={orientation}
      tabIndex={0}
      onDoubleClick={onReset}
      onKeyDown={handleKeyDown}
      onPointerDown={(event) => {
        if (event.button !== 0) return;
        lastPosition.current = coordinate(event);
        setDragging(true);
        event.currentTarget.setPointerCapture(event.pointerId);
        document.body.classList.add("is-resizing", `resize-${orientation}`);
        event.preventDefault();
      }}
      onPointerMove={(event) => {
        if (lastPosition.current === null) return;
        const next = coordinate(event);
        const delta = next - lastPosition.current;
        lastPosition.current = next;
        onDrag(delta);
      }}
      onPointerUp={finishDragging}
      onPointerCancel={finishDragging}
    >
      <span />
    </div>
  );
}
