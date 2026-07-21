import type { Metadata } from "next";
import { Plus_Jakarta_Sans, Lora, Roboto_Mono, Geist } from "next/font/google";
import "./globals.css";
import { cn } from "@/lib/utils";
import { AuthProvider } from "@/context/AuthContext";

const geist = Geist({subsets:['latin'],variable:'--font-sans'});

const fontSerif = Lora({
  subsets: ["latin"],
  variable: "--font-serif",
});

const fontMono = Roboto_Mono({
  subsets: ["latin"],
  variable: "--font-mono",
});

export const metadata: Metadata = {
  title: "PIPE - AI-Powered CI/CD Platform",
  description: "PIPE automates software delivery from code to deployment with AI-powered pull request reviews.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={cn("h-full", "font-sans", geist.variable)}>
      <body className={`${geist.variable} ${fontSerif.variable} ${fontMono.variable} antialiased min-h-full flex flex-col`}>
        <AuthProvider>
          {children}
        </AuthProvider>
      </body>
    </html>
  );
}
