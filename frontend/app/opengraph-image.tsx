import { ImageResponse } from "next/og";

export const alt = "Recode — private, account-free image tools";
export const size = { width: 1200, height: 630 };
export const contentType = "image/png";

export default function OpenGraphImage() {
  return new ImageResponse(
    (
      <div
        style={{
          alignItems: "stretch",
          background: "#f7f6f3",
          color: "#242625",
          display: "flex",
          fontFamily: "sans-serif",
          height: "100%",
          padding: "64px",
          width: "100%",
        }}
      >
        <div
          style={{
            alignItems: "stretch",
            background: "#fdfdfb",
            border: "2px solid #cfcdc6",
            borderRadius: "28px",
            display: "flex",
            flexDirection: "column",
            justifyContent: "space-between",
            padding: "56px",
            width: "100%",
          }}
        >
          <div style={{ alignItems: "center", display: "flex", gap: "20px" }}>
            <div
              style={{
                alignItems: "center",
                background: "#292b2a",
                borderRadius: "14px",
                color: "#fdfdfb",
                display: "flex",
                fontSize: "30px",
                fontWeight: 700,
                height: "68px",
                justifyContent: "center",
                letterSpacing: "-2px",
                width: "68px",
              }}
            >
              R/
            </div>
            <div style={{ fontSize: "36px", fontWeight: 700, letterSpacing: "-1.5px" }}>
              Recode
            </div>
          </div>

          <div style={{ display: "flex", flexDirection: "column", gap: "24px" }}>
            <div
              style={{
                color: "#55748d",
                fontSize: "24px",
                fontWeight: 700,
              }}
            >
              Private image workspace
            </div>
            <div
              style={{
                fontSize: "72px",
                fontWeight: 700,
                letterSpacing: "-4px",
                lineHeight: 1.02,
                maxWidth: "940px",
              }}
            >
              Image tools, without the clutter.
            </div>
            <div style={{ color: "#747773", fontSize: "27px" }}>
              14 focused tools · No account · No advertising
            </div>
          </div>

          <div style={{ display: "flex", gap: "12px" }}>
            {["Convert", "Compress", "Resize", "Crop", "Adjust"].map((label) => (
              <div
                key={label}
                style={{
                  background: "#e8eef3",
                  borderRadius: "999px",
                  color: "#3e596d",
                  display: "flex",
                  fontSize: "20px",
                  fontWeight: 600,
                  padding: "10px 20px",
                }}
              >
                {label}
              </div>
            ))}
          </div>
        </div>
      </div>
    ),
    size,
  );
}
