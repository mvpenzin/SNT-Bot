import { useState } from "react";
import { Layout } from "@/components/Layout";
import {
  useContacts,
  useCreateContact,
  useUpdateContact,
  useDeleteContact,
} from "@/hooks/use-bot";
import { insertSntContactSchema, type InsertSntContact, type SntContact } from "@shared/schema";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Textarea } from "@/components/ui/textarea";
import { Plus, Trash2, Search, Edit2, User } from "lucide-react";

export default function Contacts() {
  const { data: contacts, isLoading } = useContacts();
  const [searchTerm, setSearchTerm] = useState("");
  const [isCreateOpen, setIsCreateOpen] = useState(false);
  const [editingContact, setEditingContact] = useState<SntContact | null>(null);

  const filteredContacts =
    contacts?.filter(
      (contact) =>
        (contact.type?.toLowerCase() || "").includes(searchTerm.toLowerCase()) ||
        (contact.value?.toLowerCase() || "").includes(searchTerm.toLowerCase()) ||
        (contact.comment?.toLowerCase() || "").includes(searchTerm.toLowerCase()),
    ) || [];

  return (
    <Layout>
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 className="text-3xl font-bold tracking-tight">Контакты</h1>
          <p className="text-muted-foreground">Управление контактами СНТ</p>
        </div>
        <Button 
          onClick={() => setIsCreateOpen(true)}
          className="bg-primary text-primary-foreground hover:bg-primary/90 shadow-lg shadow-primary/20"
        >
          <Plus className="w-4 h-4 mr-2" />
          Добавить
        </Button>
      </div>

      <div className="glass-panel rounded-xl border border-border/50 overflow-hidden mt-6">
        <div className="p-4 border-b border-border/50 bg-muted/20 flex gap-4">
          <div className="relative flex-1 max-w-sm">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <Input
              placeholder="Поиск..."
              className="pl-9 bg-background/50 border-border/50 focus-visible:ring-1"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
            />
          </div>
        </div>

        <div className="overflow-x-auto">
          <Table>
            <TableHeader>
              <TableRow className="hover:bg-transparent border-border/50">
                <TableHead className="w-[80px]">Приор.</TableHead>
                <TableHead className="w-[200px]">Тип</TableHead>
                <TableHead className="min-w-[200px]">Значение</TableHead>
                <TableHead>Дополнительно</TableHead>
                <TableHead>Комментарий</TableHead>
                <TableHead className="w-[120px] text-right">Действия</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {isLoading ? (
                <TableRow>
                  <TableCell colSpan={6} className="h-24 text-center">
                    Загрузка контактов...
                  </TableCell>
                </TableRow>
              ) : filteredContacts.length === 0 ? (
                <TableRow>
                  <TableCell
                    colSpan={6}
                    className="h-48 text-center text-muted-foreground"
                  >
                    Контакты не найдены.
                  </TableCell>
                </TableRow>
              ) : (
                filteredContacts.map((contact) => (
                  <TableRow
                    key={contact.prior}
                    className="hover:bg-muted/30 border-border/50"
                  >
                    <TableCell className="font-medium text-center">
                      {contact.prior}
                    </TableCell>
                    <TableCell>
                      <div className="flex items-center gap-2">
                        <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center text-primary">
                          <User className="w-4 h-4" />
                        </div>
                        {contact.type}
                      </div>
                    </TableCell>
                    <TableCell className="text-muted-foreground">
                      {contact.value}
                    </TableCell>
                    <TableCell>
                      {contact.adds || "-"}
                    </TableCell>
                    <TableCell>
                      {contact.comment || "-"}
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex justify-end gap-2">
                        <Button
                          variant="ghost"
                          size="icon"
                          onClick={() => setEditingContact(contact)}
                          className="hover:bg-primary/10 hover:text-primary"
                        >
                          <Edit2 className="w-4 h-4" />
                        </Button>
                        <DeleteContactButton
                          prior={contact.prior}
                          type={contact.type}
                        />
                      </div>
                    </TableCell>
                  </TableRow>
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </div>

      <ContactDialog
        open={isCreateOpen}
        onOpenChange={setIsCreateOpen}
      />
      
      {editingContact && (
        <ContactDialog
          open={!!editingContact}
          onOpenChange={(open) => !open && setEditingContact(null)}
          contact={editingContact}
        />
      )}
    </Layout>
  );
}

function ContactDialog({
  open,
  onOpenChange,
  contact,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  contact?: SntContact;
}) {
  const createMutation = useCreateContact();
  const updateMutation = useUpdateContact();
  const isPending = createMutation.isPending || updateMutation.isPending;

  const form = useForm<InsertSntContact>({
    resolver: zodResolver(insertSntContactSchema),
    defaultValues: contact ? {
      type: contact.type,
      value: contact.value,
      adds: contact.adds || "",
      comment: contact.comment || "",
    } : {
      type: "",
      value: "",
      adds: "",
      comment: "",
    },
  });

  const onSubmit = (data: InsertSntContact) => {
    if (contact) {
      updateMutation.mutate({ prior: contact.prior, contact: data }, {
        onSuccess: () => {
          onOpenChange(false);
          form.reset();
        },
      });
    } else {
      createMutation.mutate(data, {
        onSuccess: () => {
          onOpenChange(false);
          form.reset();
        },
      });
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px] bg-card border-border">
        <DialogHeader>
          <DialogTitle>{contact ? "Изменить контакт" : "Добавить контакт"}</DialogTitle>
          <DialogDescription>
            {contact ? "Обновите информацию о контакте" : "Создание нового контакта СНТ"}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-4 pt-4">
          <div className="space-y-2">
            <Label htmlFor="type">Тип (например, Охрана)</Label>
            <Input
              id="type"
              {...form.register("type")}
              placeholder="Введите тип контакта"
            />
            {form.formState.errors.type && (
              <p className="text-xs text-red-500">
                {form.formState.errors.type.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="value">Значение (например, номер телефона)</Label>
            <Input
              id="value"
              {...form.register("value")}
              placeholder="Введите значение"
            />
            {form.formState.errors.value && (
              <p className="text-xs text-red-500">
                {form.formState.errors.value.message}
              </p>
            )}
          </div>
          <div className="space-y-2">
            <Label htmlFor="adds">Дополнительно</Label>
            <Input
              id="adds"
              {...form.register("adds")}
              placeholder="Дополнительная информация"
            />
          </div>
          <div className="space-y-2">
            <Label htmlFor="comment">Комментарий</Label>
            <Textarea
              id="comment"
              {...form.register("comment")}
              placeholder="Ваш комментарий"
            />
          </div>
          <DialogFooter className="pt-4">
            <Button type="submit" disabled={isPending}>
              {isPending ? "Сохранение..." : (contact ? "Сохранить" : "Создать")}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  );
}

function DeleteContactButton({ prior, type }: { prior: number; type: string }) {
  const { mutate, isPending } = useDeleteContact();

  return (
    <AlertDialog>
      <AlertDialogTrigger asChild>
        <Button
          variant="ghost"
          size="icon"
          className="hover:bg-red-500/10 hover:text-red-500"
        >
          <Trash2 className="w-4 h-4" />
        </Button>
      </AlertDialogTrigger>
      <AlertDialogContent className="bg-card border-border">
        <AlertDialogHeader>
          <AlertDialogTitle>Удалить контакт?</AlertDialogTitle>
          <AlertDialogDescription>
            Вы уверены, что хотите удалить контакт <strong>{type}</strong>? Это действие нельзя отменить.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Отмена</AlertDialogCancel>
          <AlertDialogAction
            onClick={() => mutate(prior)}
            className="bg-red-500 hover:bg-red-600 text-white"
            disabled={isPending}
          >
            {isPending ? "Удаление..." : "Удалить"}
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
