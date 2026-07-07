// Server / Client 両方から import される共通型。
// Supabase browser client 等の runtime 依存を持たないことが重要。

export class ApiError extends Error {
  readonly status: number;
  constructor(status: number, message: string) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
  get unauthorized(): boolean {
    return this.status === 401;
  }
  get notFound(): boolean {
    return this.status === 404;
  }
}

export type AuthUser = {
  id: string;
  email: string;
  name: string;
};
