export default function DashboardLoading() {
  return (
    <div className="min-h-screen bg-pulse-bg pl-[220px]">
      <div className="p-6 max-w-6xl">
        <div className="h-8 w-48 bg-pulse-surface rounded animate-pulse mb-8" />
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-8">
          {[1, 2, 3, 4].map((i) => (
            <div key={i} className="h-28 bg-pulse-surface border border-pulse-border rounded-xl animate-pulse" />
          ))}
        </div>
        <div className="h-[280px] bg-pulse-surface border border-pulse-border rounded-xl animate-pulse mb-8" />
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div className="h-48 bg-pulse-surface border border-pulse-border rounded-xl animate-pulse" />
          <div className="h-48 bg-pulse-surface border border-pulse-border rounded-xl animate-pulse" />
        </div>
      </div>
    </div>
  );
}
