#include <LiquidCrystal_I2C.h>
#include <Wire.h>
#include <HardwareSerial.h>
#include <nanoFORTH.h>
#include <string.h>

uint8_t g_x, g_y;

LiquidCrystal_I2C lcd(0x27, 16, 4); // Replace 0x27 with your I2C address

size_t pop_sane() {

  size_t r = n4_pop();

  if (r > 16)
    return 0;

  return r;
}

void setup() {
  lcd.init();      // Initialize the LCD
  lcd.backlight(); // Turn on the backlight

  Serial.begin(115200);
  while(!Serial);

  n4_setup();
  n4_api(1, forth_display);
  n4_api(2, forth_move_cursor);
  n4_api(3, forth_clear_cursor);
}

void forth_display() {
  //lcd.clear();
  lcd.setCursor(g_x, g_y);

	size_t count = pop_sane();

	// reversing the order so the string
	// is printed correctly
	for (size_t i = count; i > 0; i--) {
    
		lcd.print((char)n4_pop());
	}
}

void forth_clear_cursor() {

  size_t upto = pop_sane();

  lcd.setCursor(g_x, g_y);

  for (size_t i = g_x; i < upto; i++) {
    lcd.print(" ");
  }
}

void forth_move_cursor() {

  g_x = pop_sane();
  g_y = pop_sane();
}

void loop() {
  n4_run();
}