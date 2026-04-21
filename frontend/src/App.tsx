/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import { useState } from "react";
import { motion, AnimatePresence } from "motion/react";
import { 
  Terminal, 
  Bell, 
  User, 
  Globe, 
  Network, 
  Activity, 
  ShieldCheck, 
  Cpu, 
  BarChart3,
  Lock,
  PiggyBank,
  Cloud,
  RefreshCw,
  Shield,
  Settings,
  Rocket,
  LayoutDashboard,
  Maximize2,
  ChevronRight
} from "lucide-react";
import { 
  AreaChart, 
  Area, 
  XAxis, 
  YAxis, 
  CartesianGrid, 
  Tooltip, 
  ResponsiveContainer 
} from "recharts";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

// --- Shared Components ---

const Navbar = ({ currentPage, onNavigate }: { currentPage: string, onNavigate: (page: string) => void }) => (
  <header className="fixed top-0 z-50 w-full bg-background/80 backdrop-blur-xl border-b border-white/5">
    <div className="flex justify-between items-center w-full px-8 py-4 max-w-[1920px] mx-auto">
      <div className="flex items-center gap-8">
        <span className="font-label font-bold text-xl uppercase tracking-widest text-primary cursor-pointer" onClick={() => onNavigate("overview")}>The Kinetic Vault</span>
        <nav className="hidden md:flex gap-10 items-center">
          <button 
            className={cn(
              "font-bold font-body tracking-tight transition-all pb-1 border-b-2 text-[13px] uppercase tracking-widest",
              currentPage === "overview" ? "text-primary border-primary" : "text-on-surface-variant border-transparent hover:text-white"
            )}
            onClick={() => onNavigate("overview")}
          >
            Overview
          </button>
          <button 
            className={cn(
              "font-bold font-body tracking-tight transition-all pb-1 border-b-2 text-[13px] uppercase tracking-widest",
              currentPage === "deploy" ? "text-primary border-primary" : "text-on-surface-variant border-transparent hover:text-white"
            )}
            onClick={() => onNavigate("deploy")}
          >
            Deploy
          </button>
          <button 
            className={cn(
              "font-bold font-body tracking-tight transition-all pb-1 border-b-2 text-[13px] uppercase tracking-widest",
              currentPage === "dashboard" ? "text-primary border-primary" : "text-on-surface-variant border-transparent hover:text-white"
            )}
            onClick={() => onNavigate("dashboard")}
          >
            Dashboard
          </button>
        </nav>
      </div>
      <div className="flex items-center gap-4">
        <div className="hidden lg:flex items-center gap-2 px-3 py-1.5 rounded bg-background border border-white/10 group cursor-default">
           <Terminal size={12} className="text-on-surface-variant group-hover:text-primary transition-colors" />
           <span className="font-label text-[9px] uppercase font-bold tracking-widest text-on-surface-variant">Console Active</span>
        </div>
        <button className="flex items-center gap-2 px-4 py-2 rounded bg-primary text-on-primary font-label text-xs uppercase font-bold tracking-wider hover:brightness-110 transition-all active:scale-95 shadow-lg shadow-primary/20">
          Force Failover
        </button>
        <div className="flex items-center gap-4 ml-2">
          <Bell size={20} className="text-on-surface-variant hover:text-white cursor-pointer transition-colors" />
          <div className="w-8 h-8 rounded-full border border-white/10 flex items-center justify-center overflow-hidden cursor-pointer hover:ring-2 hover:ring-primary transition-all">
            <img 
              alt="Developer Profile" 
              className="w-full h-full object-cover" 
              src="https://lh3.googleusercontent.com/aida-public/AB6AXuDuNnei5s3Nn_QHTjqkuGmTcSzNHdsENFl47jTD8hLyanQb1p4m-91dOn6c9IbCWG6iQM2Cr_F8E2kGl2gjMfgQq1WRY7a4Vae6qMZiknetjJZ0A2WZt0tJ05mlIajwRPH1UR7MIHzlSsAza4VaHy_NmHTGiRGgXnl81Kx7gDT-_nKIPqzryPoSsopMlXCjKzYsBqtRY1DAALjkS9INXnnn7V1Z1R-FYjj-QuSQIWP3ZzE6dhXYo9wp94u-LLjQnFqvAEqSL1T5xfY"
              referrerPolicy="no-referrer"
            />
          </div>
        </div>
      </div>
    </div>
  </header>
);

// --- Overview Page Components ---

const HeroDiagram = () => (
  <div className="relative h-[600px] w-full">
    <div className="absolute inset-0 bg-primary/10 blur-[120px] rounded-full" />
    <div className="relative grid grid-cols-12 grid-rows-12 gap-4 h-full w-full">
      <motion.div 
        initial={{ opacity: 0, y: -20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.8 }}
        className="col-span-12 row-span-2 flex justify-center items-center"
      >
        <div className="glass-panel p-4 rounded-xl border border-primary/20 flex flex-col items-center gap-2 shadow-2xl">
          <div className="w-12 h-12 rounded bg-primary-container flex items-center justify-center">
            <User className="text-primary" size={24} />
          </div>
          <span className="font-label text-[10px] text-primary uppercase font-bold tracking-widest whitespace-nowrap">Global End Users</span>
        </div>
      </motion.div>

      <motion.div 
        initial={{ scale: 0.8, opacity: 0 }}
        animate={{ scale: 1, opacity: 1 }}
        transition={{ delay: 0.4, duration: 0.6 }}
        className="col-start-4 col-span-6 row-start-4 row-span-3 flex justify-center items-center relative"
      >
        <div className="absolute top-[-50px] left-1/2 w-[1px] h-[50px] bg-gradient-to-b from-primary to-transparent" />
        <div className="absolute bottom-[-50px] left-0 w-[50%] h-[1px] bg-gradient-to-r from-transparent to-primary/30" />
        <div className="absolute bottom-[-50px] right-0 w-[50%] h-[1px] bg-gradient-to-l from-transparent to-primary/30" />
        <div className="bg-surface-container-low p-6 rounded-3xl border-2 border-primary shadow-[0_0_50px_rgba(148,204,255,0.3)] flex flex-col items-center text-center z-20">
          <Cpu className="text-primary mb-2" size={32} />
          <span className="font-label text-xs text-on-surface font-bold uppercase tracking-wider">Vault DNS Core</span>
          <span className="font-label text-[9px] text-secondary">Active Monitoring</span>
        </div>
      </motion.div>

      <motion.div 
        initial={{ x: -50, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ delay: 0.6, duration: 0.8 }}
        className="col-span-5 row-start-8 row-span-4 glass-panel rounded-2xl border border-white/10 p-6 flex flex-col justify-between"
      >
        <div className="flex justify-between items-start">
          <span className="font-label text-[10px] text-on-surface-variant uppercase tracking-widest">Region A</span>
          <div className="px-2 py-0.5 rounded bg-secondary/10 text-secondary font-label text-[8px] uppercase font-bold border border-secondary/20">Primary</div>
        </div>
        <div className="flex flex-col items-center py-4">
          <img 
            alt="AWS Infrastructure" 
            className="w-16 h-16 rounded-xl mb-4 object-cover border border-primary/20 shadow-xl" 
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuBsOmdrkCYaxa82EMSkMBEEOg5bCeZjqLaMTX8n1jkzja_s1Z7xK4inoBSegioUvDhiWH2ALkyWQ68Cog5fuWFFyi1dgTqWn0IdirDXX7y4UDGVqNWNbEczZxISuKb7bWbmo0syFKnGg8JX_8CQ3TalSHPCD2GWiSEg121oEFd1Q4q-nh1bFtAl0C4GTM6yXYV8qGQYP_CgjebvUQcLRUZdVshDGUXH5HIJ59AGzFrmIrS0J1UOkGOdK886-tbKrlfrGYhbnqeL9F4"
            referrerPolicy="no-referrer"
          />
          <span className="font-label text-sm text-on-surface font-medium uppercase">AWS US-East</span>
        </div>
        <div className="h-1 w-full bg-secondary/10 rounded-full overflow-hidden">
          <motion.div 
            initial={{ width: 0 }}
            animate={{ width: "100%" }}
            transition={{ delay: 1.2, duration: 1 }}
            className="h-full bg-secondary" 
          />
        </div>
      </motion.div>

      <motion.div 
        initial={{ x: 50, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        transition={{ delay: 0.8, duration: 0.8 }}
        className="col-start-8 col-span-5 row-start-8 row-span-4 glass-panel rounded-2xl border border-white/10 p-6 flex flex-col justify-between"
      >
        <div className="flex justify-between items-start">
          <span className="font-label text-[10px] text-on-surface-variant uppercase tracking-widest">Region B</span>
          <div className="px-2 py-0.5 rounded bg-tertiary/10 text-tertiary font-label text-[8px] uppercase font-bold border border-tertiary/20">Failover</div>
        </div>
        <div className="flex flex-col items-center py-4">
          <img 
            alt="GCP Infrastructure" 
            className="w-16 h-16 rounded-xl mb-4 object-cover border border-tertiary/20 shadow-xl opacity-80" 
            src="https://lh3.googleusercontent.com/aida-public/AB6AXuBM66UD9NDCdBQg8G2gCzj5tDL4YLiELiyz4oF7dglyrnc2Z6oayZZcUuiJky8340xf_vMXOkvFGGlu15E0t2NeMWSk2lvURtCllSF5y90manZJGsqxfI70hXD1IG4CMIZzb8Z7cAWtqKJRuXZbOSSBSqG_U9RfHAuEvI0MPHjt7V06OT4IK-mcMRKUhzpgPt5zD0eVdjLlzyLFuhuFKgCbI6jluxqj_WFfZuTV5RujhPN73aQfRhM03A4Bk_Iw9O7t_MBqOGvtx1I"
            referrerPolicy="no-referrer"
          />
          <span className="font-label text-sm text-on-surface font-medium uppercase">GCP EU-West</span>
        </div>
        <div className="h-1 w-full bg-surface-container-highest rounded-full overflow-hidden">
          <motion.div 
            initial={{ width: 0 }}
            animate={{ width: "20%" }}
            transition={{ delay: 1.5, duration: 0.5 }}
            className="h-full bg-tertiary" 
          />
        </div>
      </motion.div>
    </div>
  </div>
);

const FeatureCard = ({ title, description, children, className = "", icon: Icon, iconColor = "primary" }: any) => (
  <motion.div 
    whileHover={{ y: -5 }}
    className={`bg-surface-container-high p-8 rounded-2xl border border-white/5 flex flex-col h-full relative overflow-hidden ${className}`}
  >
    {Icon && (
      <div className={`w-12 h-12 rounded flex items-center justify-center mb-6 ${
        iconColor === "primary" ? "bg-primary/10 text-primary" : 
        iconColor === "secondary" ? "bg-secondary/10 text-secondary" : 
        "bg-tertiary/10 text-tertiary"
      }`}>
        <Icon size={24} />
      </div>
    )}
    <h3 className="font-headline text-2xl font-bold text-on-surface mb-4">{title}</h3>
    <p className="text-on-surface-variant text-sm leading-relaxed mb-6">{description}</p>
    {children}
  </motion.div>
);

const OverviewPage = () => (
  <main className="pt-24">
    <section className="relative px-8 pt-20 pb-32 overflow-hidden max-w-[1920px] mx-auto">
      <div className="max-w-7xl mx-auto grid grid-cols-1 lg:grid-cols-2 gap-20 items-center">
        <motion.div 
          initial={{ opacity: 0, x: -30 }}
          animate={{ opacity: 1, x: 0 }}
          transition={{ duration: 0.8 }}
          className="z-10"
        >
          <div className="inline-flex items-center gap-2 px-3 py-1 mb-8 rounded-full bg-secondary/10 border border-secondary/20">
            <motion.span 
              animate={{ scale: [1, 1.2, 1] }} 
              transition={{ repeat: Infinity, duration: 2 }}
              className="w-2 h-2 rounded-full bg-secondary" 
            />
            <span className="font-label text-[10px] font-bold text-secondary uppercase tracking-[0.2em]">System Status: Optimal</span>
          </div>
          <h1 className="text-6xl md:text-8xl font-headline font-extrabold tracking-tighter text-on-surface leading-[0.9] mb-8">
            Architectural <span className="text-primary italic">Precision</span> for the Multicloud.
          </h1>
          <p className="text-xl text-on-surface-variant font-body max-w-xl leading-relaxed mb-12">
            The Kinetic Vault orchestrates Zero Downtime Fault Tolerance and Automatic Cost Optimization across global infrastructures. Engineered for reliability, built for scale.
          </p>
          <div className="flex flex-col sm:flex-row gap-4 mb-16">
            <button className="luminous-gradient px-10 py-5 rounded font-label font-bold text-on-primary uppercase tracking-[0.2em] shadow-2xl shadow-primary/20 hover:brightness-110 hover:scale-[1.02] transition-all active:scale-95">
              Get Started
            </button>
            <button className="bg-surface-container-highest/50 px-10 py-5 rounded font-label font-bold text-on-surface uppercase tracking-[0.2em] border border-white/10 hover:bg-surface-bright transition-all active:scale-95">
              View Documentation
            </button>
          </div>
          <div className="grid grid-cols-2 gap-12 pt-12 border-t border-white/5">
            <div>
              <span className="font-label text-4xl font-bold text-primary">0ms</span>
              <p className="font-label text-[10px] uppercase tracking-[0.3em] text-on-surface-variant mt-2">Switchover Latency</p>
            </div>
            <div>
              <span className="font-label text-4xl font-bold text-secondary">32%</span>
              <p className="font-label text-[10px] uppercase tracking-[0.3em] text-on-surface-variant mt-2">Avg Cost Reduction</p>
            </div>
          </div>
        </motion.div>
        <div className="hidden lg:block">
          <HeroDiagram />
        </div>
      </div>
    </section>

    <section className="px-8 py-32 bg-surface-container-low/30">
      <div className="max-w-7xl mx-auto">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <FeatureCard 
            className="md:col-span-2 min-h-[400px] justify-between"
            title="Zero Downtime Fault Tolerance"
            description="Our kinetic steering engine detects failures in real-time. Before your users notice a blip, traffic is re-routed across cloud boundaries with zero session loss."
          >
            <div className="mt-auto flex flex-wrap gap-4 relative z-10">
              <div className="flex items-center gap-2 bg-background px-4 py-3 rounded-lg border border-primary/20">
                <Network className="text-primary" size={16} />
                <span className="font-label text-[10px] uppercase font-bold tracking-widest">Multicloud Mesh</span>
              </div>
              <div className="flex items-center gap-2 bg-background px-4 py-3 rounded-lg border border-secondary/20">
                <ShieldCheck className="text-secondary" size={16} />
                <span className="font-label text-[10px] uppercase font-bold tracking-widest">SLO Guaranteed</span>
              </div>
            </div>
            <div className="absolute top-0 right-0 w-80 h-full opacity-5 pointer-events-none grayscale">
              <img 
                alt="Network visualization" 
                className="w-full h-full object-cover" 
                src="https://lh3.googleusercontent.com/aida-public/AB6AXuA0bmRbL_TY0u85g4LmlONvr3hXq4nvWMJzdL1SuDiHean70z9HHBjlnvoEJflZncnMmWR7XHSWX4AvO-_Uu3Uvquqas0epuah2cXd044KJRLqlReHKsPqojLrEFg9m1WNHN-fTDa1cnaYICoAsb-eUcGv0byHY8a4IcTpprDOT0gQTp0VKYiyc-38kboperzZtmjVpOBsXWgu0h5kQ-pA9VWXSEnIvgX_CcaKeCfM6Rrk4i7Tw47yw7BSFSbmdBxxwOBYTUdrtSyY"
                referrerPolicy="no-referrer"
              />
            </div>
          </FeatureCard>

          <FeatureCard 
            icon={PiggyBank}
            iconColor="tertiary"
            title="Cost Optimization"
            description="Automated logic dynamically shifts workloads to the most cost-efficient region based on spot pricing and egress fees."
          >
            <div className="mt-auto pt-8 border-t border-white/5">
              <span className="font-label text-[10px] uppercase font-bold text-tertiary tracking-[0.2em] block mb-2">Live Efficiency Score</span>
              <div className="text-4xl font-label font-bold text-on-surface">94.8%</div>
            </div>
          </FeatureCard>

          <FeatureCard 
            icon={BarChart3}
            iconColor="primary"
            title="Deep Telemetry"
            description="Every packet, every heartbeat, every failover event is logged with sub-millisecond precision. Total visibility across the stack."
          />

          <FeatureCard 
            className="md:col-span-2"
            title="Hardened Security"
            description="Encrypted transit between all nodes. Vault-level key management ensures your architecture is as secure as it is resilient."
          >
            <div className="flex items-center justify-end absolute right-8 bottom-8 pointer-events-none">
              <div className="relative w-32 h-32">
                <motion.div 
                  animate={{ rotate: 360 }}
                  transition={{ repeat: Infinity, duration: 10, ease: "linear" }}
                  className="absolute inset-0 border-[6px] border-primary/5 rounded-full" 
                />
                <motion.div 
                  animate={{ rotate: -360 }}
                  transition={{ repeat: Infinity, duration: 15, ease: "linear" }}
                  className="absolute inset-2 border-[4px] border-secondary/5 rounded-full" 
                />
                <div className="absolute inset-0 flex items-center justify-center">
                  <div className="w-16 h-16 bg-primary-container rounded-2xl flex items-center justify-center border-2 border-primary/50 shadow-[0_0_40px_rgba(148,204,255,0.4)]">
                    <Lock className="text-primary" size={28} />
                  </div>
                </div>
              </div>
            </div>
          </FeatureCard>
        </div>
      </div>
    </section>

    <section className="px-8 py-48 text-center relative overflow-hidden">
      <div className="absolute inset-0 z-0 pointer-events-none">
        <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[800px] opacity-20 bg-[radial-gradient(circle_at_center,_var(--color-primary)_0%,_transparent_70%)] blur-[100px]" />
      </div>
      <div className="max-w-4xl mx-auto relative z-10">
        <motion.h2 
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true }}
          className="text-5xl md:text-7xl font-headline font-extrabold text-on-surface mb-10 tracking-tight leading-tight"
        >
          Ready to fortify your <br />infrastructure?
        </motion.h2>
        <p className="text-xl text-on-surface-variant mb-16 max-w-2xl mx-auto">
          Join the world's most resilient engineering teams and eliminate downtime for good.
        </p>
        <motion.button 
          whileHover={{ scale: 1.05 }}
          whileTap={{ scale: 0.95 }}
          className="luminous-gradient px-12 py-6 rounded-xl font-label font-bold text-on-primary uppercase tracking-[0.3em] shadow-2xl shadow-primary/30"
        >
          Deploy Your First Vault
        </motion.button>
        <p className="mt-10 font-label text-[10px] text-on-surface-variant uppercase tracking-[0.4em] opacity-60">
          No Credit Card Required • Enterprise Ready
        </p>
      </div>
    </section>
  </main>
);

// --- Sidebar ---

const Sidebar = ({ activeItem }: { activeItem: string }) => (
  <aside className="w-64 fixed left-0 top-16 bottom-0 border-r border-white/5 bg-background flex flex-col pt-8 z-40">
    <div className="px-8 mb-10">
      <h3 className="font-label text-[10px] font-bold uppercase tracking-[0.3em] text-on-surface/40 mb-1">Operator-01</h3>
      <p className="font-label text-[10px] font-bold uppercase tracking-[0.2em] text-primary">Global Admin</p>
    </div>
    
    <nav className="flex-grow space-y-1">
      {[
        { icon: Network, label: "Infrastructure" },
        { icon: Rocket, label: "Deployments" },
        { icon: Activity, label: "Telemetry" },
        { icon: Shield, label: "Security" },
        { icon: Settings, label: "Settings" },
      ].map((item) => (
        <div 
          key={item.label}
          className={cn(
            "flex items-center gap-4 px-8 py-4 cursor-pointer transition-all border-l-4",
            activeItem === item.label 
              ? "bg-primary/5 text-primary border-primary" 
              : "text-on-surface-variant border-transparent hover:bg-white/5 hover:text-white"
          )}
        >
          <item.icon size={18} />
          <span className="font-label text-xs uppercase font-bold tracking-widest">{item.label}</span>
        </div>
      ))}
    </nav>

    <div className="p-8">
      <div className="bg-secondary/5 border border-secondary/20 p-4 rounded-lg flex flex-col gap-1">
        <div className="flex items-center gap-2">
          <motion.div 
            animate={{ opacity: [1, 0.4, 1] }} 
            transition={{ repeat: Infinity, duration: 2 }}
            className="w-2 h-2 rounded-full bg-secondary shadow-[0_0_10px_rgba(102,221,139,0.5)]" 
          />
          <span className="font-label text-[9px] font-bold uppercase tracking-widest text-secondary">System Status: Optimal</span>
        </div>
        <p className="font-label text-[7px] uppercase tracking-widest text-on-surface-variant/40 ml-4">Region: US-EAST-1</p>
      </div>
    </div>
  </aside>
);

// --- Dashboard Page ---

const chartData = [
  { time: "11:00", aws: 85, gcp: 70 },
  { time: "11:05", aws: 78, gcp: 82 },
  { time: "11:10", aws: 92, gcp: 65 },
  { time: "11:15", aws: 88, gcp: 75 },
  { time: "11:20", aws: 75, gcp: 90 },
  { time: "11:25", aws: 95, gcp: 85 },
  { time: "11:30", aws: 82, gcp: 78 },
  { time: "11:35", aws: 88, gcp: 88 },
  { time: "11:40", aws: 91, gcp: 82 },
];

const DashboardPage = () => (
  <div className="pl-64 pt-24 min-h-screen bg-surface-container-lowest">
    <Sidebar activeItem="Deployments" />
    <main className="p-10 max-w-[1600px] mx-auto">
      {/* Hero Banner */}
      <div className="relative w-full h-48 rounded-3xl overflow-hidden mb-10 border border-white/5 shadow-2xl">
        <img 
          src="https://picsum.photos/seed/datacenter/1600/400?blur=5" 
          alt="Datacenter background" 
          className="w-full h-full object-cover opacity-20"
          referrerPolicy="no-referrer"
        />
        <div className="absolute inset-0 bg-gradient-to-r from-background via-background/40 to-transparent p-12 flex flex-col justify-center">
          <h2 className="text-4xl font-headline font-extrabold text-on-surface mb-2">Vault Active Deployment</h2>
          <p className="font-label text-[10px] text-on-surface/40 uppercase font-bold tracking-[0.4em]">Global Failover Node: V-04-NORTH</p>
        </div>
      </div>

      {/* Stats Cluster */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-10">
        {[
          { label: "Provider Cluster", value: "AWS (Primary)", icon: Cloud, status: "Active/Healthy", meta: "Latency: 14ms", region: "us-east-1", color: "text-secondary" },
          { label: "Backup Cluster", value: "GCP (Backup)", icon: Cloud, status: "Standby/Healthy", meta: "Sync Lag: 0.02s", region: "europe-west1", color: "text-primary" },
          { label: "Stream Engine", value: "Data Sync", icon: RefreshCw, status: "Synchronized", meta: "Throughput: 1.2 GB/s", region: "block #98,231,012", color: "text-secondary" },
          { label: "Decision Matrix", value: "Failover Ctrl", icon: ShieldCheck, status: "Monitoring", meta: "Alert State: Nominal", region: "Auto-Repair: ENABLED", color: "text-tertiary" },
        ].map((card) => (
          <div key={card.value} className="bg-surface-container p-6 rounded-2xl border border-white/5 shadow-xl hover:border-white/10 transition-all group">
            <div className="flex justify-between items-start mb-6">
              <div>
                <p className="font-label text-[9px] text-on-surface-variant uppercase font-bold tracking-widest mb-1">{card.label}</p>
                <h4 className="text-lg font-headline font-bold text-on-surface">{card.value}</h4>
              </div>
              <card.icon size={20} className="text-on-surface-variant group-hover:text-primary transition-colors" />
            </div>
            
            <div className="inline-flex items-center px-3 py-1 bg-white/5 rounded-full mb-6 border border-white/5">
              <span className={cn("w-2 h-2 rounded-full mr-2", card.color === "text-secondary" ? "bg-secondary" : card.color === "text-primary" ? "bg-primary" : "bg-tertiary")} />
              <span className={cn("font-label text-[8px] font-bold uppercase tracking-widest", card.color)}>{card.status}</span>
            </div>

            <div className="space-y-1">
              <div className="flex justify-between text-[10px] font-label font-medium">
                <span className="text-on-surface-variant uppercase">{card.meta.split(':')[0]}</span>
                <span className="text-on-surface">{card.meta.split(':')[1]}</span>
              </div>
              <div className="flex justify-between text-[10px] font-label font-medium">
                <span className="text-on-surface-variant uppercase">{card.region.split(':')[0]}</span>
                <span className="text-on-surface">{card.region.split(':')[1] || card.region}</span>
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Heartbeat Chart */}
        <div className="lg:col-span-2 bg-surface-container p-8 rounded-3xl border border-white/5 shadow-xl">
          <div className="flex justify-between items-center mb-8">
            <div className="flex items-center gap-3">
              <div className="w-1.5 h-6 bg-primary rounded-full" />
              <h3 className="font-headline text-lg font-bold">Real-time Connection Heartbeat</h3>
            </div>
            <div className="flex items-center gap-6">
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-primary" />
                <span className="font-label text-[9px] text-on-surface-variant uppercase font-bold">AWS</span>
              </div>
              <div className="flex items-center gap-2">
                <span className="w-2 h-2 rounded-full bg-on-surface-variant" />
                <span className="font-label text-[9px] text-on-surface-variant uppercase font-bold">GCP</span>
              </div>
            </div>
          </div>

          <div className="h-[300px] w-full">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={chartData}>
                <defs>
                  <linearGradient id="colorAws" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#94ccff" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#94ccff" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorGcp" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#c6c6cd" stopOpacity={0.1}/>
                    <stop offset="95%" stopColor="#c6c6cd" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#2c3543" vertical={false} />
                <XAxis 
                  dataKey="time" 
                  axisLine={false} 
                  tickLine={false} 
                  tick={{ fill: "#909097", fontSize: 10, fontFamily: "Space Grotesk" }} 
                />
                <YAxis hide />
                <Tooltip 
                  contentStyle={{ backgroundColor: "#17202d", border: "1px solid #2c3543", borderRadius: "12px" }}
                  labelStyle={{ color: "#94ccff", marginBottom: "4px", fontSize: "12px", fontFamily: "Space Grotesk" }}
                />
                <Area type="monotone" dataKey="aws" stroke="#94ccff" strokeWidth={3} fillOpacity={1} fill="url(#colorAws)" />
                <Area type="monotone" dataKey="gcp" stroke="#909097" strokeWidth={3} fillOpacity={1} fill="url(#colorGcp)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>

          <div className="mt-8 flex justify-between items-center">
             <span className="font-label text-[9px] text-on-surface-variant/40 uppercase tracking-[0.2em]">REF_ID: KINETIC_OS_0019</span>
             <div className="flex items-center gap-2">
               <span className="font-label text-[10px] font-bold text-secondary uppercase tracking-widest">Stable: 99.998% Uptime</span>
             </div>
          </div>
        </div>

        {/* Recent Events */}
        <div className="bg-surface-container p-8 rounded-3xl border border-white/5 shadow-xl flex flex-col h-full">
          <div className="flex items-center gap-3 mb-10">
            <div className="w-1.5 h-6 bg-tertiary rounded-full" />
            <h3 className="font-headline text-lg font-bold">Recent Events</h3>
          </div>

          <div className="space-y-8 flex-grow">
            {[
              { title: "Health check passed", time: "12:04:12", meta: "Global Cluster", color: "bg-secondary" },
              { title: "Sync completed", time: "11:58:30", meta: "Storage Layer", color: "bg-secondary" },
              { title: "Backup initialized", time: "11:45:01", meta: "GCP-EU-WEST", color: "bg-primary" },
              { title: "Threshold warning suppressed", time: "11:30:15", meta: "Latency Monitor", color: "bg-tertiary" },
            ].map((event) => (
              <div key={event.time} className="flex gap-4">
                <div className={cn("w-2 h-2 rounded-full mt-1.5 shrink-0", event.color)} />
                <div>
                  <h4 className="font-body text-sm font-bold text-on-surface mb-0.5">{event.title}</h4>
                  <p className="font-label text-[10px] text-on-surface-variant uppercase tracking-widest leading-relaxed">
                    {event.time} • {event.meta}
                  </p>
                </div>
              </div>
            ))}
          </div>

          <button className="w-full mt-12 py-4 rounded-xl font-label text-[10px] uppercase font-bold tracking-[0.2em] border border-white/5 hover:bg-white/5 transition-all">
            View All Logs
          </button>
        </div>
      </div>
    </main>
  </div>
);

// --- Deploy Page ---

const StrategyCard = ({ title, description, active = false, icon: Icon, onClick }: any) => (
  <div 
    onClick={onClick}
    className={cn(
      "p-6 rounded-xl border transition-all cursor-pointer flex gap-4",
      active 
        ? "bg-primary/5 border-primary shadow-[0_0_20px_rgba(148,204,255,0.1)]" 
        : "bg-surface-container border-white/5 hover:border-white/10"
    )}
  >
    <div className={cn(
      "w-4 h-4 rounded-full border-2 mt-1 flex items-center justify-center shrink-0",
      active ? "border-primary" : "border-on-surface-variant"
    )}>
      {active && <div className="w-2 h-2 rounded-full bg-primary" />}
    </div>
    <div>
      <div className="flex items-center gap-2 mb-1">
        <h4 className="font-headline font-bold text-sm uppercase tracking-wider">{title}</h4>
        {Icon && <Icon size={14} className="text-primary hover:text-white transition-colors" />}
      </div>
      <p className="font-label text-[10px] text-on-surface-variant uppercase leading-tight tracking-widest">{description}</p>
    </div>
  </div>
);

const predictiveData = [
  { h: "00:00", v: 40 },
  { h: "01:00", v: 30 },
  { h: "02:00", v: 50 },
  { h: "03:00", v: 80 },
  { h: "NOW", v: 100, active: true },
  { h: "+1H", v: 70 },
  { h: "+2H", v: 90 },
  { h: "+3H", v: 45 },
  { h: "+4H", v: 30 },
];

const DeployPage = () => {
  const [strategy, setStrategy] = useState("ha");
  const [replicas, setReplicas] = useState(32);

  return (
    <div className="pl-64 pt-24 min-h-screen bg-surface-container-lowest">
      <Sidebar activeItem="Deployments" />
      <main className="p-10 max-w-[1600px] mx-auto">
        <header className="mb-12">
          <p className="font-label text-[10px] text-on-surface-variant uppercase font-bold tracking-[0.4em] mb-2">Vault Protocol / Initializer</p>
          <h2 className="text-6xl font-headline font-extrabold text-on-surface tracking-tight">Initialize Deployment</h2>
          <p className="text-on-surface-variant max-w-2xl mt-4 font-body leading-relaxed">
            Configure high-performance clusters with architectural precision. Our failover engine ensures 99.999% availability across disparate cloud providers.
          </p>
        </header>

        <div className="flex flex-col lg:flex-row gap-10">
          <div className="flex-grow space-y-8">
            {/* Primary Config */}
            <div className="bg-surface-container p-10 rounded-3xl border border-white/5 shadow-xl relative overflow-hidden h-fit">
               <div className="absolute top-8 right-8 text-on-surface/5">
                 <Network size={80} />
               </div>
               <h3 className="font-headline text-[10px] font-bold uppercase tracking-[0.4em] mb-10 text-on-surface-variant">Primary Configuration</h3>
               
               <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-10">
                 <div>
                   <label className="block font-label text-[10px] uppercase font-bold tracking-widest text-on-surface-variant mb-4">Application Name</label>
                   <input 
                     type="text" 
                     defaultValue="vault-api-production"
                     className="w-full bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface focus:outline-none focus:border-primary transition-colors"
                   />
                 </div>
                 <div>
                   <label className="block font-label text-[10px] uppercase font-bold tracking-widest text-on-surface-variant mb-4">Container Image</label>
                   <div className="relative">
                     <input 
                       type="text" 
                       defaultValue="ghcr.io/kinetic/vault-api:latest"
                       className="w-full bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface focus:outline-none focus:border-primary transition-colors"
                     />
                     <div className="absolute right-4 top-1/2 -translate-y-1/2 text-on-surface-variant/40">
                       <Rocket size={14} />
                     </div>
                   </div>
                 </div>
               </div>

               <div className="grid grid-cols-1 md:grid-cols-12 gap-10 items-end">
                 <div className="md:col-span-4">
                    <label className="block font-label text-[10px] uppercase font-bold tracking-widest text-on-surface-variant mb-4">Ports</label>
                    <input 
                       type="text" 
                       defaultValue="80, 443, 8080"
                       className="w-full bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface focus:outline-none focus:border-primary transition-colors"
                    />
                 </div>
                 <div className="md:col-span-8">
                    <div className="flex justify-between items-center mb-6">
                      <label className="font-label text-[10px] uppercase font-bold tracking-widest text-on-surface-variant">Replicas</label>
                      <span className="bg-background border border-white/10 px-4 py-1.5 rounded-lg font-label text-[12px] font-bold text-primary">{replicas}</span>
                    </div>
                    <input 
                      type="range" 
                      min="1" 
                      max="128" 
                      value={replicas}
                      onChange={(e) => setReplicas(parseInt(e.target.value))}
                      className="w-full h-1.5 bg-background rounded-full appearance-none cursor-pointer accent-primary"
                    />
                 </div>
               </div>
            </div>

            {/* Env Variables */}
            <div className="bg-surface-container p-10 rounded-3xl border border-white/5 shadow-xl h-fit">
              <div className="flex justify-between items-center mb-10">
                <h3 className="font-headline text-[10px] font-bold uppercase tracking-[0.4em] text-on-surface-variant">Environment Variables</h3>
                <button className="flex items-center gap-2 font-label text-[10px] uppercase text-primary font-bold tracking-widest hover:text-white transition-colors">
                  <div className="text-xl">+</div> Add Variable
                </button>
              </div>

              <div className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                   <div className="bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface-variant uppercase">DB_CONNECTION</div>
                   <div className="bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface">vault://psql-cluster-01</div>
                </div>
                <div className="grid grid-cols-2 gap-4">
                   <input placeholder="KEY" className="bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface-variant focus:outline-none focus:border-primary transition-colors" />
                   <input placeholder="VALUE" className="bg-background border border-white/10 rounded-lg px-6 py-4 font-label text-[11px] text-on-surface focus:outline-none focus:border-primary transition-colors" />
                </div>
              </div>
            </div>
          </div>

          <div className="w-full lg:w-96 space-y-8">
            {/* Strategy */}
            <div className="bg-surface-container p-8 rounded-3xl border border-white/5 shadow-2xl flex flex-col h-fit">
              <h3 className="font-headline text-[10px] font-bold uppercase tracking-[0.4em] text-on-surface-variant mb-8 text-center">Deployment Strategy</h3>
              
              <div className="space-y-4">
                <StrategyCard 
                  title="Cost-Optimized" 
                  description="Leverage spot instances & regional routing" 
                  active={strategy === "cost"}
                  icon={Bell}
                  onClick={() => setStrategy("cost")}
                />
                <StrategyCard 
                  title="High-Availability" 
                  description="Multi-cloud mesh with instant failover" 
                  active={strategy === "ha"}
                  icon={ShieldCheck}
                  onClick={() => setStrategy("ha")}
                />
                <StrategyCard 
                  title="Custom Build" 
                  description="Manual resource allocation & weights" 
                  active={strategy === "custom"}
                  icon={Settings}
                  onClick={() => setStrategy("custom")}
                />
              </div>

              <div className="mt-12 pt-12 border-t border-white/5">
                 <div className="flex justify-between items-center mb-8 px-2">
                   <span className="font-label text-[10px] text-on-surface-variant uppercase tracking-[0.3em]">Est. Cost/Hour</span>
                   <span className="font-label text-xl font-bold text-secondary">$0.42</span>
                 </div>
                 <button className="w-full luminous-gradient py-6 rounded-xl font-headline font-bold text-[12px] uppercase tracking-[0.3em] shadow-2xl shadow-primary/30 hover:brightness-110 active:scale-95 transition-all">
                   Deploy to Vault
                 </button>
                 <p className="mt-6 text-center font-label text-[8px] text-on-surface-variant/40 uppercase tracking-widest">
                   Transaction will be verified via <br /> AUTH-01 protocol
                 </p>
              </div>
            </div>

            {/* Predictive Chart */}
            <div className="bg-surface-container p-8 rounded-3xl border border-white/5 shadow-2xl">
               <div className="flex items-center gap-3 mb-10">
                 <div className="w-1.5 h-6 bg-primary rounded-full" />
                 <h3 className="font-headline text-lg font-bold">Predictive Load</h3>
               </div>
               
               <div className="flex items-end justify-between h-32 px-2 gap-2">
                 {predictiveData.map((d, i) => (
                   <div key={i} className="flex flex-col items-center gap-3 w-full">
                     <motion.div 
                        initial={{ height: 0 }}
                        animate={{ height: `${d.v}%` }}
                        transition={{ delay: i * 0.1, duration: 0.8 }}
                        className={cn(
                          "w-full rounded-t-sm",
                          d.active ? "bg-primary shadow-[0_0_15px_rgba(148,204,255,0.4)]" : "bg-white/10"
                        )} 
                     />
                     <span className={cn("font-label text-[8px] uppercase font-bold", d.active ? "text-primary" : "text-on-surface-variant/40")}>{d.h}</span>
                   </div>
                 ))}
               </div>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};

// --- Main App Component ---

export default function App() {
  const [currentPage, setCurrentPage] = useState("overview");

  return (
    <div className="bg-background text-on-surface font-body selection:bg-primary/30 selection:text-primary min-h-screen">
      <Navbar currentPage={currentPage} onNavigate={setCurrentPage} />
      
      <AnimatePresence mode="wait">
        <motion.div
          key={currentPage}
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -10 }}
          transition={{ duration: 0.3 }}
        >
          {currentPage === "overview" && <OverviewPage />}
          {currentPage === "deploy" && <DeployPage />}
          {currentPage === "dashboard" && <DashboardPage />}
        </motion.div>
      </AnimatePresence>
    </div>
  );
}

