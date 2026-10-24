import { Trash2 } from "lucide-react";
import { Modal } from "./Modal";
import { Button } from "./Button";

interface Props {
  title: string;
  message: string;
  busy?: boolean;
  onConfirm: () => void;
  onClose: () => void;
}

export function ConfirmDeleteModal({ title, message, busy, onConfirm, onClose }: Props) {
  return (
    <Modal title={title} onClose={onClose}>
      <div className="confirm-body">
        <div className="confirm-icon">
          <Trash2 size={23} />
        </div>
        <h3>هل أنت متأكد من الحذف؟</h3>
        <p>{message}</p>
        <div className="modal-actions">
          <Button variant="secondary" onClick={onClose}>
            إلغاء
          </Button>
          <Button variant="danger" disabled={busy} onClick={onConfirm}>
            {busy ? "جاري الحذف..." : "تأكيد الحذف"}
          </Button>
        </div>
      </div>
    </Modal>
  );
}
