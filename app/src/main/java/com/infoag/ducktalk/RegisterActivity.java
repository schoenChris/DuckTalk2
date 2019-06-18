package com.infoag.ducktalk;

import android.content.Intent;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class RegisterActivity extends AppCompatActivity {

    private EditText usernameEdit;
    private EditText passwordEdit;
    private EditText rePasswordEdit;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_register);

        usernameEdit = findViewById(R.id.usernameField);
        passwordEdit = findViewById(R.id.passwordField);
        rePasswordEdit = findViewById(R.id.rePasswordField);

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

        passwordEdit.requestFocus();
    }

    public void attemptRegister() {

        // get input from user
        String usernameText = usernameEdit.getText().toString();
        String passwordText = passwordEdit.getText().toString();
        String rePasswordText = rePasswordEdit.getText().toString();

        // compare passwords, continue when equal
        if (passwordText.equals(rePasswordText)) {

            // continue

        }

    }

}
