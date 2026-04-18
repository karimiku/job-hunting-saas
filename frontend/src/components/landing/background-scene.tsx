"use client";

const PHOTO_URL =
  "https://images.unsplash.com/photo-1464822759023-fed622ff2c3b?auto=format&fit=crop&w=2400&q=80";

export function BackgroundScene() {
  return (
    <div
      aria-hidden
      style={{
        position: "fixed",
        inset: 0,
        overflow: "hidden",
        pointerEvents: "none",
        zIndex: 0,
      }}
    >
      {/* Sunrise gradient base */}
      <div
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(180deg, #E8E0FA 0%, #FBE4ED 22%, #FFE0D1 42%, #FFEDD6 58%, #EAF2FC 85%, #E0EAFA 100%)",
        }}
      />

      {/* Sea-of-clouds photo, tonally reduced */}
      <div
        style={{
          position: "absolute",
          inset: 0,
          backgroundImage: `url("${PHOTO_URL}")`,
          backgroundSize: "cover",
          backgroundPosition: "center center",
          opacity: 0.55,
          filter: "saturate(1.1) contrast(0.92) brightness(1.05)",
        }}
      />

      {/* Color-unify overlay */}
      <div
        style={{
          position: "absolute",
          inset: 0,
          background:
            "linear-gradient(180deg, rgba(232,224,250,0.35) 0%, rgba(251,228,237,0.2) 30%, rgba(255,237,214,0.18) 55%, rgba(234,242,252,0.3) 85%, rgba(224,234,250,0.45) 100%)",
        }}
      />

      {/* Aurora blobs */}
      <div
        style={{
          position: "absolute",
          top: "-10%",
          left: "-10%",
          width: "55vw",
          height: "55vw",
          borderRadius: "50%",
          background:
            "radial-gradient(circle, rgba(154,123,232,0.3), transparent 65%)",
          filter: "blur(70px)",
          opacity: 0.7,
        }}
      />
      <div
        style={{
          position: "absolute",
          top: "10%",
          right: "-8%",
          width: "45vw",
          height: "45vw",
          borderRadius: "50%",
          background:
            "radial-gradient(circle, rgba(255,179,152,0.35), transparent 65%)",
          filter: "blur(60px)",
          opacity: 0.7,
        }}
      />
      <div
        style={{
          position: "absolute",
          bottom: "-15%",
          left: "20%",
          width: "35vw",
          height: "35vw",
          borderRadius: "50%",
          background:
            "radial-gradient(circle, rgba(168,216,245,0.45), transparent 65%)",
          filter: "blur(70px)",
          opacity: 0.7,
        }}
      />

      {/* Rising light particles — hope */}
      {Array.from({ length: 24 }).map((_, i) => {
        const left = (i * 53 + 7) % 100;
        const delay = (i * 0.31) % 6;
        const dur = 10 + (i % 4) * 2;
        const size = 3 + (i % 3) * 2;
        const hue = ["#FFD94A", "#FFB398", "#FFFFFF", "#A8D8F5"][i % 4];
        return (
          <div
            key={i}
            style={{
              position: "absolute",
              left: `${left}%`,
              bottom: "-20px",
              width: size,
              height: size,
              borderRadius: "50%",
              background: hue,
              boxShadow: `0 0 ${size * 3}px ${hue}, 0 0 ${size * 6}px ${hue}88`,
              animation: `lpRise ${dur}s linear ${delay}s infinite`,
              opacity: 0,
            }}
          />
        );
      })}

      {/* Shimmer light beam */}
      <div
        style={{
          position: "absolute",
          top: "30%",
          left: 0,
          width: "40%",
          height: "2px",
          background:
            "linear-gradient(90deg, transparent, rgba(255,255,255,0.7), transparent)",
          filter: "blur(2px)",
          animation: "lpShimmer 10s ease-in-out infinite",
        }}
      />
    </div>
  );
}
