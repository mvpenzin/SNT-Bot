import { pgTable, text, serial, timestamp, boolean } from "drizzle-orm/pg-core";
import { createInsertSchema } from "drizzle-zod";
import { comment } from "postcss";
import { number, z } from "zod";

// --- Users Table (from Telegram) ---
export const sntUsers = pgTable("snt_users", {
  id: serial("id").primaryKey(),
  created: timestamp("created").defaultNow(),
  modified: timestamp("modified").defaultNow(),
  user_id: text("user_id").notNull().unique(),
  user_name: text("user_name"),
  user_fio: text("first_name"),
  lastName: text("last_name"),
});

// --- Contacts Table ---
export const sntContacts = pgTable("snt_contacts", {
  id: serial("id").primaryKey(),
  created: timestamp("created").defaultNow(),
  modified: timestamp("modified").defaultNow(),
  prior: serial("prior").notNull().unique(),
  type: text("type").notNull(),
  value: text("value").notNull(),
  adds: text("adds"),
  comment: text("comment"),
});

// --- Bot Logs (for frontend display) ---
export const botLogs = pgTable("bot_logs", {
  id: serial("id").primaryKey(),
  created: timestamp("created").defaultNow(),
  modified: timestamp("modified").defaultNow(),
  level: text("level").notNull(), // "INFO", "WARN", "ERROR"
  message: text("message").notNull(),
  details: text("details"),
});

// --- Schemas ---
export const insertSntUserSchema = createInsertSchema(sntUsers).omit({
  id: true,
  created: true,
});
export const insertSntContactSchema = createInsertSchema(sntContacts).omit({
  id: true,
  created: true,
});
export const insertBotLogSchema = createInsertSchema(botLogs).omit({
  id: true,
  created: true,
});

// --- Types ---
export type SntUser = typeof sntUsers.$inferSelect;
export type InsertSntUser = z.infer<typeof insertSntUserSchema>;

export type SntContact = typeof sntContacts.$inferSelect;
export type InsertSntContact = z.infer<typeof insertSntContactSchema>;

export type BotLog = typeof botLogs.$inferSelect;
export type InsertBotLog = z.infer<typeof insertBotLogSchema>;
