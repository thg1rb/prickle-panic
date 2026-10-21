import { cardAttributes } from "../game/cocoCardVariant";

export const th = {
  brand: "ปิกเกิล แพนิก!",
  picnic: "ปิกนิกเบอร์รีของโคโค่",
  pageTitle: "ปิกนิกเบอร์รีของโคโค่",
  pageDescription:
    "เกมไพ่จับคู่ของโคโค่ เลือกไพ่ให้ต่างจากการ์ดอ้างอิงไม่เกินหนึ่งอย่างก่อนหมดเวลา",
  footer: "เล่นในเครื่องนี้",
  coco: "โคโค่",
  times: "ครั้ง",
  seconds: "วินาที",
  cardAttributes: {
    bodyCream: "ตัวสีครีม",
    bodyBrown: "ตัวสีน้ำตาล",
    longQuills: "ขนยาว",
    shortQuills: "ขนสั้น",
    glasses: "ใส่แว่น",
    noGlasses: "ไม่ใส่แว่น",
    scarf: "มีผ้าพันคอ",
    noScarf: "ไม่มีผ้าพันคอ",
    strawberry: "มีสตรอว์เบอร์รี",
    noStrawberry: "ไม่มีสตรอว์เบอร์รี",
  },
  tutorial: {
    label: "วิธีเล่น",
    title: "มองให้ไว แล้วแตะให้ถูก",
  },
  home: {
    eyebrow: "ปิกนิกเบอร์รีของโคโค่",
    title: "พร้อมลุยหรือยัง?",
    description:
      "หาไพ่ที่หน้าตาเหมือนหรือแตกต่างจากโคโค่ไม่เกินหนึ่งอย่าง ก่อนหมดเวลา!",
    start: "เริ่มเล่นคนเดียว",
    mode: "โหมดฝึกเล่น",
    duration: "รอบละ 60 วินาที",
  },
  game: {
    back: "กลับหน้าหลัก",
    reference: "การ์ดอ้างอิง",
    score: "คะแนน",
    combo: "ต่อเนื่อง",
    timer: "เวลาที่เหลือ",
    instruction: "เลือกการ์ดที่ต่างจากการ์ดอ้างอิงไม่เกินหนึ่งอย่าง",
    correct: "เก่งมาก! โคโค่เปลี่ยนการ์ดอ้างอิงให้แล้ว",
    incorrect: "ยังไม่ใช่นะ ลองสังเกตโคโค่อีกครั้ง",
    cooldown: "พักนิ้วสักครู่นะ",
    expired: "หมดเวลาแล้ว!",
    cardPrefix: "การ์ดโคโค่",
    board: "กระดานการ์ดโคโค่",
  },
  result: {
    title: "หมดเวลาแล้ว!",
    score: "คะแนนที่ทำได้",
    replay: "เล่นอีกครั้ง",
    home: "กลับหน้าหลัก",
  },
  guest: {
    loading: "กำลังตรวจสอบเซสชันของคุณ...",
    recoveryTitle: "ยังเชื่อมต่อไม่ได้",
    title: "ยินดีต้อนรับสู่ปิกนิกของโคโค่",
    description: "ตั้งชื่อเล่นก่อน แล้วมาเล่นกับโคโค่กันนะ",
    greeting: "สวัสดี",
    usernameLabel: "ชื่อผู้เล่น",
    usernameHint: "ใช้ชื่อ 2–16 ตัวอักษร และไม่เกิน 256 ไบต์",
    placeholder: "พิมพ์ชื่อที่อยากให้เพื่อนเห็น",
    submit: "เริ่มปิกนิก",
    creating: "กำลังเตรียมเซสชัน...",
    retry: "ลองตรวจสอบอีกครั้ง",
    wait: "ลองใหม่ใน",
    seconds: "วินาที",
    graphemeCount: "กรุณาใช้ชื่อ 2–16 ตัวอักษร",
    tooLong: "ชื่อนี้ยาวเกินไป กรุณาใช้ไม่เกิน 256 ไบต์",
    errors: {
      usernameInvalid: "ชื่อนี้ยังใช้ไม่ได้ ลองตรวจสอบตัวอักษรอีกครั้งนะ",
      requestInvalid: "ข้อมูลชื่อไม่ถูกต้อง กรุณาลองใหม่",
      requestTooLarge: "ข้อมูลชื่อยาวเกินไป กรุณาใช้ชื่อที่สั้นลง",
      originNotAllowed: "ไม่สามารถส่งคำขอจากหน้านี้ได้ กรุณาโหลดเกมใหม่",
      rateLimited: "มีการเริ่มเซสชันหลายครั้งเกินไป กรุณารอสักครู่",
      sessionRequired: "เซสชันหมดอายุแล้ว กรุณาใส่ชื่ออีกครั้ง",
      sessionUnavailable: "บริการยังไม่พร้อม กรุณาลองอีกครั้ง",
      unexpected: "เกิดข้อผิดพลาดชั่วคราว กรุณาลองอีกครั้ง",
    },
  },
} as const;

export function describeCardTh(variant: number): string {
  const [bodyCream, longQuills, glasses, scarf, strawberry] =
    cardAttributes(variant);
  return [
    bodyCream ? th.cardAttributes.bodyCream : th.cardAttributes.bodyBrown,
    longQuills ? th.cardAttributes.longQuills : th.cardAttributes.shortQuills,
    glasses ? th.cardAttributes.glasses : th.cardAttributes.noGlasses,
    scarf ? th.cardAttributes.scarf : th.cardAttributes.noScarf,
    strawberry ? th.cardAttributes.strawberry : th.cardAttributes.noStrawberry,
  ].join(" ");
}
