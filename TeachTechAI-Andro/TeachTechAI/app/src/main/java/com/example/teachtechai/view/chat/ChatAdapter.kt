package com.example.teachtechai.view.chat

import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.TextView
import androidx.recyclerview.widget.RecyclerView
import com.example.teachtechai.R
import com.example.teachtechai.data.Message
import com.example.teachtechai.data.response.ChatResponse

class ChatAdapter(private val messages: MutableList<Message>) :
    RecyclerView.Adapter<RecyclerView.ViewHolder>() {

    companion object{
        private const val VIEW_TYPE_USER = 1
        private const val VIEW_TYPE_AI = 2
    }
    override fun getItemViewType(position: Int): Int {
        return if (messages[position].isUser) VIEW_TYPE_USER else VIEW_TYPE_AI
    }
    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): RecyclerView.ViewHolder {
        return if (viewType == VIEW_TYPE_USER) {
            val view = LayoutInflater.from(parent.context).inflate(R.layout.message_item_user, parent, false)
            UserMessageViewHolder(view)
        } else{
            val view = LayoutInflater.from(parent.context).inflate(R.layout.message_item_ai, parent, false)
            AiMessageViewHolder(view)
        }
    }

    override fun onBindViewHolder(holder: RecyclerView.ViewHolder, position: Int) {
        val message = messages[position]
        if(holder is UserMessageViewHolder){
            holder.messageTextView.text = message.text
        }else if(holder is AiMessageViewHolder){
            holder.messageTextView.text = message.text
        }
    }
    class AiMessageViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        val messageTextView : TextView = itemView.findViewById(R.id.tvMessage)
    }

    class UserMessageViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        val messageTextView : TextView = itemView.findViewById(R.id.tvMessage)
    }

    override fun getItemCount(): Int {
        return messages.size
    }



}