import Link from "next/link";

export function Header() {
  return (
    <header className="w-full h-16 border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60 flex items-center px-6 shrink-0 z-50">
      <Link
        href="/"
        className="font-bold text-xl hover:opacity-80 transition-opacity"
      >
        ToDo App
      </Link>
    </header>
  );
}
