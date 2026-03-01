import Navigation from "@/components/Navigation";
import Hero from "@/components/Hero";
import SocialProof from "@/components/SocialProof";
import ProblemSection from "@/components/ProblemSection";
import HowItWorks from "@/components/HowItWorks";
import DashboardPreview from "@/components/DashboardPreview";
import Constellation from "@/components/Constellation";
import PrivacyCTA from "@/components/PrivacyCTA";
import Footer from "@/components/Footer";

export default function LandingPage() {
  return (
    <main className="min-h-screen bg-pulse-bg selection:bg-pulse-primary selection:text-white">
      <Navigation />
      <Hero />
      <SocialProof />
      <ProblemSection />
      <HowItWorks />
      <DashboardPreview />
      <Constellation />
      <PrivacyCTA />
      <Footer />
    </main>
  );
}
