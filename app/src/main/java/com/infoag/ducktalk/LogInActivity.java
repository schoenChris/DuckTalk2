package com.infoag.ducktalk;

import android.content.Intent;
import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.view.View;
import android.widget.Button;
import android.widget.EditText;

public class LogInActivity extends AppCompatActivity {

    private EditText usernameEdit;
    private EditText passwordEdit;

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_log_in);

        usernameEdit = findViewById(R.id.usernameField);
        passwordEdit = findViewById(R.id.passwordField);

        Button logInButton = findViewById(R.id.logInButton);
        logInButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                attemptLogin();
            }
        });

        Button registerButton = findViewById(R.id.registerButton);
        registerButton.setOnClickListener(new View.OnClickListener() {
            @Override
            public void onClick(View v) {
                String usernameText = usernameEdit.getText().toString();
                Intent intent = new Intent(LogInActivity.this, RegisterActivity.class);
                intent.putExtra("username_message", usernameText);
                startActivity(intent);
            }
        });
    }

    public void attemptLogin() {

        // Missing


    }

}
