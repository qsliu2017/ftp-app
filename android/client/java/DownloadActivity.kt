package com.example.ftpclient

import android.content.Intent
import android.os.Bundle
import android.view.View
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import fm.Fm
import kotlinx.android.synthetic.main.activity_download.*

class DownloadActivity : AppCompatActivity(), View.OnClickListener {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_download)
        downloadBtn.setOnClickListener(this)
        Fm.setFileManager(
            Connection.FileManagerImpl(
                ContextCompat.getExternalFilesDirs(
                    this,
                    null
                )[0]
            )
        )
    }

    override fun onClick(v: View?) {
        when (v?.id) {
            R.id.downloadBtn -> {
                val local = localName.text.toString().trim()
                val remote = remoteName.text.toString().trim()

                // Download a file
                try {
                    Connection.getCon()?.retrieve(local, remote)
                    val from = intent.getStringExtra("from").toString()
                    val dialog = AlertDialog.Builder(this).setMessage("Download successfully!")
                    if (from == "main") {
                        dialog.setPositiveButton(
                            "OK"
                        ) { _, _ ->
                            startActivity(Intent(this, MainActivity::class.java))
                        }
                    } else {
                        dialog.setPositiveButton(
                            "OK"
                        ) { _, _ ->
                            startActivity(Intent(this, UserActivity::class.java))
                        }
                    }
                    dialog.create().show()
                } catch (e: Exception) {
                    val error = Connection.exceptionHandle(e)
                    AlertDialog.Builder(this).setMessage(error)
                        .setPositiveButton("OK", null).create().show()
                }
            }
        }
    }
}