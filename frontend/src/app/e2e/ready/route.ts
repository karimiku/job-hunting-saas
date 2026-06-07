import { NextResponse } from "next/server";

export function GET() {
  if (process.env.PLAYWRIGHT_E2E_AUTH !== "true") {
    return NextResponse.json({ message: "Not found" }, { status: 404 });
  }

  return NextResponse.json({ ok: true });
}
