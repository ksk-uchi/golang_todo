import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Header } from "@/app/components/Header";
import { Footer } from "@/app/components/Footer";
import { Toaster } from "@/app/components/ui/sonner";
import { ThemeProvider } from "@/app/components/theme-provider";
import { Providers } from "@/app/providers";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "ToDo App",
  description: "You can manage your todos",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body
        className={`${geistSans.variable} ${geistMono.variable} antialiased h-screen flex flex-col overflow-hidden`}
      >
        <Providers>
          <ThemeProvider
            attribute="class"
            defaultTheme="dark"
            enableSystem={false}
            disableTransitionOnChange
          >
            <Header />
            <main className="flex-1 overflow-y-auto flex flex-col">
              <div className="container mx-auto px-4 py-8 max-w-2xl flex-1 flex flex-col">
                {children}
              </div>
            </main>
            <Footer />
            <Toaster richColors />
          </ThemeProvider>
        </Providers>
      </body>
    </html>
  );
}
