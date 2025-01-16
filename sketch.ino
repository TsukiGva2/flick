//#include <EnableInterrupt.h>
#include <LiquidCrystal_I2C.h>
//#include <Wire.h>
//#include <HardwareSerial.h>
#include <nanoFORTH.h>

#define LABEL_COUNT 11

const char* labels[] = {
  "PORTAL   My",
  "ATLETAS  ",
  "REGIST.  ",
  "COMUNICANDO ",
  "LEITOR ",
  "LTE/4G: ",
  "WIFI: ",
  "IP: ",
  "LOCAL: ",
  "PROVA: ",
  "PING: "
};
const int labels_len[LABEL_COUNT] = {
  11,9,9,12,7,8,6,4,7,7,6
};

#define VALUE_COUNT 4

const char* values[] = {
  "WEB",
  "CONECTAD",
  "DESLIGAD",
  "AUTOMATIC",
  "OK",
  "X"
};

const char code[] PROGMEM =          ///< define preload Forth code here
": lbl 5 API ;\n"
": fwd 2 API ;\n"
": lit API fwd ;\n"
": num 4 lit ;\n"
": val 6 lit ;\n"
": ip  7 lit ;\n"
": atn 1 lit ;\n";

uint8_t g_x, g_y;

LiquidCrystal_I2C lcd(0x27, 16, 4); // Replace 0x27 with your I2C address

void setup() {
  lcd.init();      // Initialize the LCD
  lcd.backlight(); // Turn on the backlight

  Serial.begin(115200);
  while(!Serial);

  n4_setup(code);
  n4_api(1, forth_antenna);
  n4_api(2, forth_fwd);
  // n4_api(3, forth_clear_cursor);

  n4_api(5, forth_label);
  n4_api(6, forth_value);
  n4_api(7, forth_ip);
  n4_api(4, forth_number);

  pinMode(7, INPUT_PULLUP);
}

void forth_value() {

  int v = n4_pop();
  if (v > VALUE_COUNT) {
    lcd.print("---");
    return;
  }

  lcd.print(values[v]);
}

void print_forthNumber() {

  int mag = n4_pop();
  int v = n4_pop();

  lcd.print(v);

  if (mag == 0)
    return;

  if (mag >= 3 && mag < 6) {

    lcd.print('K');

    return;
  }

  lcd.print('M');
}

void forth_antenna() {

  forth_clear_line(4);
  lcd.setCursor(0, g_y);
  lcd.print("A1: ");
  print_forthNumber();
  lcd.print("  A2: ");
  print_forthNumber();
  forth_fwd();
  forth_clear_line(4);
  lcd.setCursor(0, g_y);
  lcd.print("A3: ");
  print_forthNumber();
  lcd.print("  A4: ");
  print_forthNumber();
}

void forth_ip() {

  lcd.print(n4_pop());
  lcd.print('.');
  lcd.print(n4_pop());
  lcd.print('.');
  lcd.print(n4_pop());
  lcd.print('.');
  lcd.print(n4_pop());
}

void forth_number() {

  lcd.print(n4_pop());
}

void forth_label() {

  static int current_labels[4] = {-1,-1,-1,-1};

  int v = n4_pop();
  if (v >= LABEL_COUNT) {
    lcd.print("-----");
    return;
  }

  if (v != current_labels[g_y]) {
    forth_clear_line(0);
    current_labels[g_y] = v;
    lcd.print(labels[v]);
  } else {
    forth_clear_line(labels_len[v]);
  }
}

void forth_clear_line(int x) {
  lcd.setCursor(x, g_y);

  for (size_t i = x; i <= 16; i++) {
    lcd.print(' ');
  }

  lcd.setCursor(x, g_y);
}

void forth_fwd() {
  ++g_y %= 4;
}

void loop() {
  n4_run();
}