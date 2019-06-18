package com.infoag.ducktalk;

import android.content.Intent;
import android.content.SharedPreferences;
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

        // get username and password from user input
        String usernameText = usernameEdit.getText().toString();
        String passwordText = passwordEdit.getText().toString();

        // create login task
        // new ServerTask(USER_VALIDATION).execute(usernameText, passwordText);


        //if (success) {

            // save login data in Shared Preferences
            SharedPreferences sp = getSharedPreferences("logIn", MODE_PRIVATE);
            SharedPreferences.Editor editor = sp.edit();
            editor.putString("username", usernameText);
            editor.putString("password", passwordText);
            editor.apply();

        //}


    }

}
