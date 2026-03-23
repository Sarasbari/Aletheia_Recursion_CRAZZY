import { Inter, Space_Grotesk, JetBrains_Mono, Caveat } from "next/font/google";
import "./globals.css";
import "./pages.css";
import "./wallet.css";
import Web3Provider from "./components/Web3Provider";

const inter = Inter({
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700", "800", "900"],
  variable: "--font-body",
});

const spaceGrotesk = Space_Grotesk({
  subsets: ["latin"],
  weight: ["300", "400", "500", "600", "700"],
  variable: "--font-brand",
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-mono",
});

const caveat = Caveat({
  subsets: ["latin"],
  weight: ["400", "500", "600", "700"],
  variable: "--font-hand",
});

export const metadata = {
  title: "ALETHEIA — Proof of Origin. On-chain.",
  description:
    "ALETHEIA anchors image authenticity on-chain. We hash it, sign it, put it on-chain. Nobody can argue with math.",
};

export default function RootLayout({ children }) {
  return (
    <html
      lang="en"
      className={`${inter.variable} ${spaceGrotesk.variable} ${jetbrainsMono.variable} ${caveat.variable}`}
    >
      <body>
        <Web3Provider>{children}</Web3Provider>
      </body>
    </html>
  );
}
