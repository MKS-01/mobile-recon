import Navbar from '@/components/Navbar';
import HeroSection from '@/components/HeroSection';
import FeaturesSection from '@/components/FeaturesSection';
import UsageSection from '@/components/UsageSection';
import InstallSection from '@/components/InstallSection';
import Footer from '@/components/Footer';

export default function Home() {
  return (
    <main className="min-h-screen">
      <Navbar />
      <HeroSection />
      <FeaturesSection />
      <UsageSection />
      <InstallSection />
      <Footer />
    </main>
  );
}
