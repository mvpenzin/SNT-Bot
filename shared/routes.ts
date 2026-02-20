import { z } from "zod";
import {
  insertSntContactSchema,
  insertBotLogSchema,
  sntUsers,
  sntContacts,
  botLogs,
} from "./schema";

export const errorSchemas = {
  validation: z.object({
    message: z.string(),
    field: z.string().optional(),
  }),
  notFound: z.object({
    message: z.string(),
  }),
  internal: z.object({
    message: z.string(),
  }),
};

export const api = {
  // --- Bot Status & Logs ---
  status: {
    get: {
      method: "GET" as const,
      path: "/api/status" as const,
      responses: {
        200: z.object({
          status: z.enum(["running", "stopped", "error"]),
          uptime: z.number(),
          lastCheck: z.string(),
        }),
      },
    },
  },
  logs: {
    list: {
      method: "GET" as const,
      path: "/api/logs" as const,
      responses: {
        200: z.array(z.custom<typeof botLogs.$inferSelect>()),
      },
    },
  },

  // --- Contacts Management (Admin) ---
  contacts: {
    list: {
      method: "GET" as const,
      path: "/api/contacts" as const,
      responses: {
        200: z.array(z.custom<typeof sntContacts.$inferSelect>()),
      },
    },
    create: {
      method: "POST" as const,
      path: "/api/contacts" as const,
      input: insertSntContactSchema,
      responses: {
        201: z.custom<typeof sntContacts.$inferSelect>(),
        400: errorSchemas.validation,
      },
    },
    update: {
      method: "PUT" as const,
      path: "/api/contacts/:prior" as const,
      parameters: {
        prior: z.coerce.number(),
      },
      input: insertSntContactSchema,
      responses: {
        200: z.custom<typeof sntContacts.$inferSelect>(),
        400: errorSchemas.validation,
        404: errorSchemas.notFound,
      },
    },
    delete: {
      method: "DELETE" as const,
      path: "/api/contacts/:prior" as const,
      parameters: {
        prior: z.coerce.number(),
      },
      responses: {
        204: z.void(),
        404: errorSchemas.notFound,
      },
    },
  },
};

export function buildUrl(
  path: string,
  params?: Record<string, string | number>,
): string {
  let url = path;
  if (params) {
    Object.entries(params).forEach(([key, value]) => {
      if (url.includes(`:${key}`)) {
        url = url.replace(`:${key}`, String(value));
      }
    });
  }
  return url;
}
