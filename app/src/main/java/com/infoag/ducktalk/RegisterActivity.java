package com.infoag.ducktalk;

import android.content.Intent;
import android.content.SharedPreferences;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;
import android.widget.TextView;

public class RegisterActivity extends AppCompatActivity {

    private EditText usernameEdit;
    private EditText passwordEdit;
    private EditText rePasswordEdit;
    private TextView errorTextView;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register);

        usernameEdit = findViewById(R.id.usernameField);
        passwordEdit = findViewById(R.id.passwordField);
        rePasswordEdit = findViewById(R.id.rePasswordField);
        errorTextView = findViewById(R.id.errorTextView);
        errorTextView.setVisibility(View.INVISIBLE);

        Button registerButton = findViewById(R.id.registerButton);
        registerButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                attemptRegister();
            }
        });

        Intent intent = getIntent();
        String usernameText = intent.getStringExtra("username_message");
        usernameEdit.setText(usernameText);
    }

    public void attemptRegister() {

        // get user input
        String username = usernameEdit.getText().toString();
        String password = passwordEdit.getText().toString();
        String rePassword = rePasswordEdit.getText().toString();

        if (username.isEmpty()) {
            setErrorText("Bitte Benutzernamen eingeben!");

        } else if (password.isEmpty()) {
            setErrorText("Bitte Passwort eingeben!");

        } else if (password.length() <= 4) {
            setErrorText("Passwort ist zu kurz!");

        } else if (!password.equals(rePassword)) {
            setErrorText("Passwörter sind nicht gleich!");

        } else {

            // register new user

            // save login data in Shared Preferences
            SharedPreferences sp = getSharedPreferences("logIn", MODE_PRIVATE);
            SharedPreferences.Editor editor = sp.edit();
            editor.putString("username", username);
            editor.putString("password", password);
            editor.apply();

            Intent intent = new Intent(this, ContactViewActivity.class);
            startActivity(intent);

        }

    }

    private void setErrorText(String errorText) {
        errorTextView.setText(errorText);
        if (errorTextView.getVisibility() != View.VISIBLE) {
            errorTextView.setVisibility(View.VISIBLE);
        }
    }

}
