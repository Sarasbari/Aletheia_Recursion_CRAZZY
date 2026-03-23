import Navbar from "./components/Navbar";
import Hero from "./components/Hero";
import HowItWorks from "./components/HowItWorks";
import WhyAletheia from "./components/WhyAletheia";
import UseCases from "./components/UseCases";
import Pricing from "./components/Pricing";
import Footer from "./components/Footer";

export default function Home() {
  return (
    <>
      <Navbar />
      <Hero />
      <hr className="section-divider" />
      <HowItWorks />
      <hr className="section-divider" />
      <WhyAletheia />
      <hr className="section-divider" />
      <UseCases />
      <hr className="section-divider" />
      <Pricing />
      <Footer />
    </>
  );
}
